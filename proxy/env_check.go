package proxy

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AddProxy takes a user defined URL and routes tests through it
func AddProxy(rawProxy string, cmd *cobra.Command) *url.URL {
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
		CheckProxy(rawProxy)
	}
	return proxyURL
}

// CheckProxy verifies the user defined or auto-detected proxy is viable for reaching public sites
func CheckProxy(rawProxy string) {
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
	log.Info("Connection OK.  Able to reach www.saucelabs.com.", resp)
}

// CheckForEnvProxies double-checks that common environment variables aren't set
func CheckForEnvProxies() {
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
