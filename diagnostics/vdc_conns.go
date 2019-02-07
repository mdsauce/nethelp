package diagnostics

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

// SauceServices sends HTTP requests to Sauce endpoints
func SauceServices(sauceEndpoints []string) {
	for _, endpoint := range sauceEndpoints {
		u, err := url.ParseRequestURI(endpoint)
		if err != nil {
			log.WithFields(log.Fields{
				"err":      err,
				"endpoint": endpoint,
			}).Debug("Could not parse endpoint.")
			fmt.Printf("[ ] %s is not reachable. Err: %v\n", endpoint, err)
			continue
		}
		log.WithFields(log.Fields{
			"url":    u,
			"IsAbs?": u.IsAbs(),
			"scheme": u.Scheme,
			"port":   u.Port,
			"path":   u.Path,
		}).Debug("URL after Parsing")

		log.Debug("Sending GET req to ", u)
		resp, err := http.Get(u.String())
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"endpoint": u,
			}).Fatalf("[ ] %s not reachable\n", u)
		}

		respOutput(resp, endpoint)
	}
}

// VdcAPI connects to the rest endpoint with env variable credentials
func VdcAPI(vdcRESTEndpoints []string) {
	username := os.Getenv("SAUCE_USERNAME")
	apiKey := os.Getenv("SAUCE_ACCESS_KEY")
	for _, endpoint := range vdcRESTEndpoints {
		log.Debug("Sending GET req to ", endpoint)
		var jsonBody = []byte(`{}`)
		req, err := http.NewRequest("GET", endpoint, bytes.NewBuffer(jsonBody))
		req.SetBasicAuth(username, apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatalf("[ ] %s not reachable\n", endpoint)
		}
		respOutput(resp, endpoint)
	}
}
