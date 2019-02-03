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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var proxy string

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
		VerboseMode(cmd)
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is $HOME/.nethelp.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print all logging levels")
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "p", "", "upstream proxy for nethelp to use.  Port should be added like my.proxy:8080")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
