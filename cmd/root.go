package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mdsauce/nethelp/connections"
	"github.com/mdsauce/nethelp/endpoints"
	"github.com/mdsauce/nethelp/proxy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string
var userProxy string

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
connections. by sending requests to services used 
during a Sauce Labs session (RDC or VDC) .`,
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
		proxyURL := proxy.AddProxy(userProxy, cmd)
		log.Info("Proxy URL: ", proxyURL)
		proxy.CheckForEnvProxies()

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
		// refine data from cli and assemble
		// endpoints/services to be tested
		whichCloud = strings.ToLower(whichCloud)
		whichDC = strings.ToLower(whichDC)
		vdcTest := endpoints.NewVDCTest(whichDC)
		headlessTest := endpoints.NewHeadlessTest(whichDC)
		rdcTest := endpoints.NewRDCTest(whichDC)
		defPublic := endpoints.NewPublicTest()
		vdcAPITest := endpoints.AssembleVDCEndpoints(whichDC)
		headlessAPITest := endpoints.AssembleHeadlessEndpoints(whichDC)

		if whichDC != "all" {
			validateDC(whichDC)
		}
		// Run the diagnostics that the user passed in
		if whichCloud != "all" {
			validateCloud(whichCloud)
			// VDC
			if whichCloud == "vdc" {
				connections.VDCServices(vdcTest.Endpoints)
				if vdcAPITest != nil {
					connections.VdcAPI(vdcAPITest.Endpoints)
				}
			}
			// RDC
			if whichCloud == "rdc" {
				connections.RDCServices(rdcTest.Endpoints)
			}
			// Headless
			if whichCloud == "headless" {
				connections.HeadlessServices(headlessTest.Endpoints)
				if headlessAPITest != nil {
					connections.HeadlessAPI(headlessAPITest.Endpoints)
				}
			}
		}

		if runTCP {
			defTCP := endpoints.NewTCPTest()
			connections.TCPConns(defTCP.Sitelist, proxyURL)
		} else if whichCloud == "all" {
			connections.VDCServices(vdcTest.Endpoints)
			connections.RDCServices(rdcTest.Endpoints)
			connections.HeadlessServices(headlessTest.Endpoints)
			if whichDC == "all" {
				connections.PublicSites(defPublic.Sitelist)
			}
			if vdcAPITest != nil {
				connections.VdcAPI(vdcAPITest.Endpoints)
			}
			if headlessAPITest != nil {
				connections.HeadlessAPI(headlessAPITest.Endpoints)
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is $HOME/.nethelp.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "print all logging levels")
	rootCmd.PersistentFlags().StringVarP(&userProxy, "proxy", "p", "", "upstream proxy for nethelp to use. Enter like -p protocol://username:password@host:port")

	rootCmd.Flags().BoolP("lucky", "l", false, "disable the proxy check at startup and instead test the proxy during execution.")
	rootCmd.Flags().Bool("tcp", false, "run TCP tests. Will always run against all endpoints.")
	rootCmd.Flags().Bool("log", false, "enables logging and creates a nethelp.log file.  Will automatically append data to the file in a non-destructive manner.")
	rootCmd.Flags().String("cloud", "all", "options are: VDC, RDC, or HEADLESS.  Select which services you'd like to test, Virtual Device Cloud, Real Device Cloud, or the Headless Cloud.")
	rootCmd.Flags().String("dc", "all", "options are: EU, NA, or EAST.  Choose which data centers you want run diagnostics against, Europe, North America(West), or North America(East).")

	// http client settings
	http.DefaultTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
}

func validateCloud(whichCloud string) {
	if whichCloud != "vdc" && whichCloud != "rdc" && whichCloud != "headless" {
		log.Fatal("The parameter is not valid.  Only 'all', 'vdc', 'rdc', or 'headless' are allowed")
	}
}

func validateDC(whichDC string) {
	if whichDC != "na" && whichDC != "eu" && whichDC != "east" {
		log.Fatal("The parameter is not valid.  Only 'all', 'na', 'eu', or 'east' are allowed")
	}
}
