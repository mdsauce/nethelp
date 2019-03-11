// Copyright © 2019 Max Dobeck
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mdsauce/nethelp/diagnostics"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var userProxy string
var publicSites, tcplist, vdcNA, vdcEU, rdcNA, rdcEU []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nethelp",
	Short: "Helper to troubleshoot problems running tests on Sauce Labs.",
	Long: `
 ___  __ _ _   _  ___ ___   / / __   ___| |_| |__   ___| |_ __  
/ __|/ _  | | | |/ __/ _ \ / / '_ \ / _ \ __| '_ \ / _ \ | '_ \ 
\__ \ (_| | |_| | (_|  __// /| | | |  __/ |_| | | |  __/ | |_) |
|___/\__,_|\__,_|\___\___/_/ |_| |_|\___|\__|_| |_|\___|_| .__/ 
                                                         |_|  
Nethelp will help find out what is blocking outbound 
connections by sending requests to 
services used during typical Sauce Labs usage.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Logging and Verbosity setup
		log.SetOutput(os.Stdout)
		log.SetLevel(log.WarnLevel)
		enableVerbose, err := cmd.PersistentFlags().GetBool("verbose")
		if err != nil {
			log.Fatal("Verbose flag broke.", err)
		}
		if enableVerbose == true {
			log.SetLevel(log.TraceLevel)
		}
		logging, err := cmd.Flags().GetBool("log")
		if err != nil {
			log.Fatal("Could not get output flag.")
		}
		if logging == true {
			log.SetFormatter(&log.TextFormatter{
				DisableColors: true,
			})
			fp, err := os.OpenFile("./nethelp.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				log.Fatal(err)
			}
			defer fp.Close()
			log.SetOutput(fp)
			log.Warn("This log only captures output from the --verbose flag.")
		}

		// Proxy setup and configuration
		proxyURL := addProxy(userProxy, cmd)
		log.Info("Proxy URL: ", proxyURL)
		checkForEnvProxies()

		// Default APIs and service endpoints
		tcplist = []string{"ondemand.saucelabs.com:443", "ondemand.saucelabs.com:80", "ondemand.saucelabs.com:8080", "ondemand.eu-central-1.saucelabs.com:80", "ondemand.eu-central-1.saucelabs.com:443", "us1.appium.testobject.com:443", "eu1.appium.testobject.com:443", "us1.appium.testobject.com:80", "eu1.appium.testobject.com:80"}
		publicSites = []string{"https://status.us-west-1.saucelabs.com","http://status.eu-central-1.saucelabs.com/", "https://www.duckduckgo.com"}
		vdcNA = []string{"https://ondemand.saucelabs.com:443", "http://ondemand.saucelabs.com:80"}
		vdcEU = []string{"http://ondemand.eu-central-1.saucelabs.com:80", "https://ondemand.eu-central-1.saucelabs.com:443"}
		rdcNA = []string{"https://us1.appium.testobject.com/wd/hub/session"}
		rdcEU = []string{"https://eu1.appium.testobject.com/wd/hub/session"}
		naVDCApi, euVDCApi := assembleVDCEndpoints()

		// Collect the flags to decide which diagnostics to run
		runTCP, err := cmd.Flags().GetBool("tcp")
		if err != nil {
			log.Fatal("Could not get the TCP flag. ", err)
		}
		whichDC, err := cmd.Flags().GetString("dc")
		if err != nil {
			log.Fatal("Could not get the dc flag. ", err)
		}
		whichCloud, err := cmd.Flags().GetString("cloud")
		if err != nil {
			log.Fatal("Could not get the cloud flag. ", err)
		}
		whichCloud = strings.ToLower(whichCloud)
		whichDC = strings.ToLower(whichDC)

		// Run the diagnostics that the user passed in
		if whichCloud != "all" {
			if whichCloud != "vdc" && whichCloud != "rdc" {
				log.Fatal("The parameter is not valid.  Only 'all', 'vdc', or 'rdc' are allowed. ", whichCloud)
			}
			// VDC and a specific region
			if whichCloud == "vdc" && whichDC != "all" {
				if whichDC == "eu" {
					diagnostics.VDCServices(vdcEU)
					diagnostics.VdcAPI(euVDCApi)
				} else if whichDC == "na" {
					diagnostics.VDCServices(vdcNA)
					diagnostics.VdcAPI(naVDCApi)
				}
			}
			// RDC and a specific region
			if whichCloud == "rdc" && whichDC != "all" {
				if whichDC == "eu" {
					diagnostics.RDCServices(rdcEU)
				} else if whichDC == "na" {
					diagnostics.RDCServices(rdcNA)
				}
			}
			// VDC and all regions OR RDC and all regions
			if whichCloud == "vdc" && whichDC == "all" {
				diagnostics.VDCServices(vdcNA)
				diagnostics.VDCServices(vdcEU)
			} else if whichCloud == "rdc" && whichDC == "all" {
				diagnostics.RDCServices(rdcEU)
				diagnostics.RDCServices(rdcNA)
			}
		}
		// Specific region and all clouds
		if whichCloud == "all" && whichDC != "all" {
			if whichDC == "eu" {
				diagnostics.VDCServices(vdcEU)
				diagnostics.RDCServices(rdcEU)
				diagnostics.VdcAPI(euVDCApi)	
			} else if whichDC == "na" {
				diagnostics.VDCServices(vdcNA)
				diagnostics.RDCServices(rdcNA)
				diagnostics.VdcAPI(naVDCApi)			
			}
		}
		if runTCP {
			diagnostics.TCPConns(tcplist, proxyURL)
		}
		if runDefault(runTCP) && whichDC == "all" && whichCloud == "all" {
			diagnostics.VDCServices(vdcNA)
			diagnostics.VDCServices(vdcEU)
			diagnostics.PublicSites(publicSites)
			diagnostics.RDCServices(rdcEU)
			diagnostics.RDCServices(rdcNA)
			diagnostics.VdcAPI(naVDCApi)
			diagnostics.VdcAPI(euVDCApi)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	log.SetFormatter(&log.TextFormatter{})
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is $HOME/.nethelp.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "print all logging levels")
	rootCmd.PersistentFlags().StringVarP(&userProxy, "proxy", "p", "", "upstream proxy for nethelp to use. Enter like -p protocol://username:password@host:port")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("lucky", "l", false, "disable the proxy check at startup and instead test the proxy during execution.")
	rootCmd.Flags().Bool("tcp", false, "run TCP tests. Will always run against all endpoints.")
	rootCmd.Flags().Bool("log", false, "enables logging and creates a nethelp.log file.  Will automatically append data to the file in a non-destructive manner.")
	rootCmd.Flags().String("cloud", "all", "options are: VDC or RDC.  Select which services you'd like to test, Virtual Device Cloud or Real Device Cloud respectively.")
	rootCmd.Flags().String("dc", "all", "options are: EU or NA.  Choose which data centers you want run diagnostics against, Europe or North America respectively.")

	// http client settings
	http.DefaultTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		log.Debug("Config file: ", cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Failed to find homedir")
		}
		home = home + "/.config"
		// Search config in home directory with name ".nethelp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nethelp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func addProxy(rawProxy string, cmd *cobra.Command) *url.URL {
	var proxyURL *url.URL
	var err error
	var disableCheck bool
	disableCheck, err = cmd.Flags().GetBool("lucky")
	if err != nil {
		log.Fatal("Something went terribly wrong disabling the check with --lucky", err)
	}

	if rawProxy != "" {
		proxyURL, err = url.Parse(rawProxy)
		if err != nil {
			log.Fatalf("Panic while setting proxy %s.  Proxy not set and program exiting. %v", rawProxy, err)
		}
		// This takes care of HTTP calls globally
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	// check that there are no env vars defining a proxy and everything works
	if disableCheck != true {
		checkProxy(rawProxy)
	}
	return proxyURL
}

func checkProxy(rawProxy string) {
	resp, err := http.Get("https://www.saucelabs.com")
	if err != nil {
		if rawProxy != "" {
			log.WithFields(log.Fields{
				"error": err,
				"msg":   "www.saucelabs.com not reachable with this proxy",
			}).Fatalf("Something is wrong with the user specified proxy %s.  It cannot be used.", rawProxy)
		} else {
			log.WithFields(log.Fields{
				"error": err,
				"msg":   "www.saucelabs.com not reachable.",
			}).Warn("You may have no internet access or a proxy may be in use.")
		}
	}
	log.Info("Proxy OK.  Able to reach www.saucelabs.com.", resp)
}

func runDefault(runTCP bool) bool {
	if runTCP {
		log.Debug("Specific test flag used.  Not running default test set.")
		return false
	}
	return true
}

// assembleVDCEndpoints interpolates user variables like
// username and their sauce api key to create a valid URI
// Returns NA api list then EU api list
func assembleVDCEndpoints() ([]string, []string) {
	if os.Getenv("SAUCE_USERNAME") == "" {
		log.Info("No Environment Variables found.  Not running VDC REST endpoint tests.")
		return nil, nil
	}
	naVDCREST := []string{""}
	euVDCREST := []string{""}
	endpoint := fmt.Sprintf("https://saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))
	naVDCREST[0] = endpoint
	endpoint = fmt.Sprintf("https://eu-central-1.saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))
	euVDCREST[0] = endpoint
	return naVDCREST, euVDCREST
}

func checkForEnvProxies() {
	proxyList := []string{"HTTP_PROXY", "HTTPS_PROXY", "PROXY", "ALL_PROXY"}
	for _, proxy := range proxyList {
		if os.Getenv(proxy) != "" {
			log.WithFields(log.Fields{
				"env var":         proxy,
				"potential proxy": os.Getenv(proxy),
			}).Warn("An environment variable for a Proxy may exist.  This will NOT be automatically used.  You must use the --proxy flag.")
		}
	}
}
