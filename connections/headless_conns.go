package connections

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

// HeadlessServices sends HTTP requests to Headless Sauce endpoints to prove
// tests could theoretically be created and the data centers are reachable
func HeadlessServices(sauceEndpoints []string) {
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
			fmt.Printf("[ ] %s not reachable\n", u)
			log.WithFields(log.Fields{
				"error":    err,
				"endpoint": u,
			}).Infof("[ ] %s not reachable\n", u)
		}

		if err == nil {
			respOutput(resp, endpoint)
		}
	}
}

// HeadlessAPI connects to Headless (us-east-1) REST endpoints to make sure
// 1) credentials work
// 2) api is reachable
// 3) api retrieves the expected data if 1 & 2 are true
func HeadlessAPI(vdcRESTEndpoints []string) {
	log.Debug("Sending out HTTP reqs to these endpoints: ", vdcRESTEndpoints)
	username := os.Getenv("SAUCE_USERNAME")
	apiKey := os.Getenv("HEADLESS_ACCESS_KEY")
	for _, endpoint := range vdcRESTEndpoints {
		log.Debug("Sending req to ", endpoint)
		var jsonBody = []byte(`{}`)
		req, err := http.NewRequest("GET", endpoint, bytes.NewBuffer(jsonBody))
		req.SetBasicAuth(username, apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[ ] %s not reachable\n", endpoint)
			log.WithFields(log.Fields{
				"error": err,
			}).Infof("[ ] %s not reachable\n", endpoint)
		}

		if err == nil {
			respOutput(resp, endpoint)
		}
	}
}
