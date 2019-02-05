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
	"net/http"
	"net/url"
	"os"

	"github.com/mdsauce/nethelp/diagnostics"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var userProxy string
var sitelist, tcplist, vdcEndpoints, rdcEndpoints []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nethelp",
	Short: "Helper to troubleshoot problems running tests on Sauce Labs.",
	Long: `Nethelp will assist with finding out what is blocking outbound 
connections from the machine by sending HTTP and TCP connections to 
services used by Sauce Labs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.WarnLevel)
		VerboseMode(cmd)
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
		}
		// log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		proxyURL := addProxy(userProxy, cmd)
		log.Info("Proxy URL: ", proxyURL)
		checkForEnvProxies()

		//default sitelists
		tcplist = []string{"ondemand.saucelabs.com:443", "ondemand.saucelabs.com:80", "ondemand.saucelabs.com:8080", "ondemand.eu-central-1.saucelabs.com:80", "ondemand.eu-central-1.saucelabs.com:443", "us1.appium.testobject.com:443", "eu1.appium.testobject.com:443", "us1.appium.testobject.com:80", "eu1.appium.testobject.com:80"}
		sitelist = []string{"https://status.saucelabs.com", "https://www.duckduckgo.com"}
		vdcEndpoints = []string{"https://ondemand.saucelabs.com:443", "http://ondemand.saucelabs.com:80", "ondemand.eu-central-1.saucelabs.com:80", "ondemand.eu-central-1.saucelabs.com:443"}
		// TODO
		rdcEndpoints = []string{"https://us1.appium.testobject.com/wd/hub/session", "https://eu1.appium.testobject.com/wd/hub/session"}
		vdcRESTEndpoints := assembleVDCEndpoints()

		runHTTP, err := cmd.Flags().GetBool("http")
		if err != nil {
			log.Fatal("Could not get the HTTP flag. ", err)
		}
		runTCP, err := cmd.Flags().GetBool("tcp")
		if err != nil {
			log.Fatal("Could not get the TCP flag. ", err)
		}
		runAPI, err := cmd.Flags().GetBool("api")
		if err != nil {
			log.Fatal("Could not get the API flag. ", err)
		}
		if runHTTP {
			diagnostics.PublicSites(sitelist)
			diagnostics.SauceServices(vdcEndpoints)
			diagnostics.RDCServices(rdcEndpoints)
		}
		if runTCP {
			diagnostics.TCPConns(tcplist, proxyURL)
		}
		if runAPI {
			diagnostics.VDCREST(vdcRESTEndpoints)
		}
		if runDefault(runHTTP, runTCP, runAPI) {
			diagnostics.PublicSites(sitelist)
			diagnostics.SauceServices(vdcEndpoints)
			diagnostics.RDCServices(rdcEndpoints)
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
	rootCmd.Flags().Bool("http", false, "run HTTP tests. Default is to run all tests.")
	rootCmd.Flags().Bool("tcp", false, "run TCP tests. Default is to only run HTTP tests.")
	// rootCmd.Flags().StringP("out", "o", time.Now().Format("20060102150405"), "optional output file for logging. Defaults to timestamp file in the current dir.  Only use if you want a custom log name.")
	rootCmd.Flags().Bool("log", false, "enables logging to the file specified by the --out flag.")
	rootCmd.Flags().Bool("api", false, "run API tests.  Requires that you have $SAUCE_USERNAME and $SAUCE_ACCESS_KEY environment variables.")
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

// VerboseMode enables Trace level logging and shows all logs
func VerboseMode(cmd *cobra.Command) {
	enableVerbose, err := cmd.PersistentFlags().GetBool("verbose")
	if err != nil {
		log.Fatal("Verbose flag broke.", err)
	}
	if enableVerbose == true {
		log.SetLevel(log.TraceLevel)
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
				"msg":   "www.saucelabs.com not reachable with this proxy",
				"resp":  *resp,
			}).Fatal("Something is wrong with the proxy specified in your environment variables.  It cannot be used.")
		}
	}
	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"response": *resp,
		}).Fatal("Something is wrong with the proxy.  It cannot was not able to reach a public website.")
	}
	log.Info("Proxy OK.  Able to reach www.saucelabs.com.", resp)
}

func runDefault(runHTTP bool, runTCP bool, runAPI bool) bool {
	if runHTTP || runTCP || runAPI {
		log.Debug("Specific test flag used.  Not running default test set.")
		return false
	}
	return true
}

func assembleVDCEndpoints() []string {
	if os.Getenv("SAUCE_USERNAME") == "" {
		log.Info("No Environment Variables found.  Not running VDC REST endpoint tests.")
		return nil
	}
	vdcRESTEndpoints := []string{""}
	endpoint := fmt.Sprintf("https://saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))
	vdcRESTEndpoints[0] = endpoint
	endpoint := fmt.Sprintf("https://eu-central-1.saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))
	vdcRESTEndpoints[1] = endpoint
	return vdcRESTEndpoints
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
