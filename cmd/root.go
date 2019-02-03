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
var sitelist, tcplist, vdcEndpoints, rdcEndpoints, vdcRESTEndpoints []string

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
		log.SetLevel(log.WarnLevel)
		VerboseMode(cmd)
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		proxyURL := addProxy(userProxy, cmd)
		log.Info("Proxy URL: ", proxyURL)

		//default sitelists
		tcplist = []string{"ondemand.saucelabs.com:443", "ondemand.saucelabs.com:80", "ondemand.saucelabs.com:8080", "us1.appium.testobject.com:443", "eu1.appium.testobject.com:443", "us1.appium.testobject.com:80", "eu1.appium.testobject.com:80"}
		sitelist = []string{"https://status.saucelabs.com", "https://www.duckduckgo.com"}
		vdcEndpoints = []string{"https://ondemand.saucelabs.com:443", "http://ondemand.saucelabs.com:80"}
		// TODO
		rdcEndpoints = []string{"https://us1.appium.testobject.com/wd/hub/session", "https://eu1.appium.testobject.com/wd/hub/session"}
		vdcRESTEndpoints = []string{"https://saucelabs.com/rest/v1/USERNAME/tunnels"}

		if runDefault(cmd) {
			diagnostics.PublicSites(sitelist)
			diagnostics.SauceServices(vdcEndpoints)
			diagnostics.RDCServices(rdcEndpoints)
			diagnostics.VDCREST(vdcRESTEndpoints)
		} else {
			runHTTP, err := cmd.Flags().GetBool("http")
			if err != nil {
				log.Fatal("Could not get the HTTP flag. ", err)
			}
			runTCP, err := cmd.Flags().GetBool("tcp")
			if err != nil {
				log.Fatal("Could not get the TCP flag. ", err)
			}
			if runHTTP {
				diagnostics.PublicSites(sitelist)
				diagnostics.SauceServices(vdcEndpoints)
				diagnostics.RDCServices(rdcEndpoints)
				diagnostics.VDCREST(vdcRESTEndpoints)
			}
			if runTCP {
				diagnostics.TCPConns(tcplist, proxyURL)
			}
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is $HOME/.nethelp.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "print all logging levels")
	rootCmd.PersistentFlags().StringVarP(&userProxy, "proxy", "p", "", "upstream proxy for nethelp to use.  Port should be added like my.proxy:8080")
	// TODO
	// root.rootCmd.PersistentFlags().StringVarP(&proxyAuth, "auth", "a", "", "authentication for upstream proxy.  use like -a username:password.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("lucky", "l", false, "feeling lucky?  Disable the proxy check at startup and find out if it works during runtime.")
	rootCmd.Flags().Bool("http", false, "run HTTP tests. Default is to run all tests.")
	rootCmd.Flags().Bool("tcp", false, "run TCP tests. Default is to run all tests.")
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
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
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
			}).Fatalf("Something is wrong with the proxy %s.  It cannot be used.", rawProxy)
		} else {
			log.WithFields(log.Fields{
				"error": err,
				"msg":   "www.saucelabs.com not reachable with this proxy",
			}).Fatal("Something is wrong with the proxy specified in your environment variables.  It cannot be used.")
		}
	}
	log.Info("Proxy OK.  Able to reach www.saucelabs.com.", resp)
}

func runDefault(cmd *cobra.Command) bool {
	runHTTP, err := cmd.Flags().GetBool("http")
	if err != nil {
		log.Fatal("Could not get the HTTP flag. ", err)
	}
	runTCP, err := cmd.Flags().GetBool("tcp")
	if err != nil {
		log.Fatal("Could not get the TCP flag. ", err)
	}
	log.Debug("Checking commands to see if we don't want to run default test set. ", runHTTP, runTCP)
	if runHTTP || runTCP {
		log.Debug("HTTP or TCP flag used.  Not running default test set.", runHTTP, runTCP)
		return false
	}
	return true
}
