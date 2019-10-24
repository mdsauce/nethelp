package connections

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// RDCServices makes connections to the main RDC endpoints to prove
// that the endpoints are reachable from the machine
func RDCServices(rdcEndpoints []string) {
	for _, endpoint := range rdcEndpoints {
		log.Debug("Sending POST req to ", endpoint)
		var jsonBody = []byte(`{"test":"this will result in an HTTP 500 resp or 401 resp."}`)
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
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

// RdcAPI connects to RDC REST endpoints to make sure
// 1) credentials work
// 2) api is reachable
// 3) api retrieves the expected data
func RdcAPI(rdcRESTEndpoints []string) {
	log.Debug("Sending out HTTP reqs to these endpoints: ", rdcRESTEndpoints)
	username := os.Getenv("RDC_USERNAME")
	apiKey := os.Getenv("RDC_ACCESS_KEY")
	for _, endpoint := range rdcRESTEndpoints {
		log.Debug("Sending GET req to ", endpoint)
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
