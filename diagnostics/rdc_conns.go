package diagnostics

import (
	"bytes"
	"net/http"

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
			log.WithFields(log.Fields{
				"error": err,
			}).Warnf("[ ] %s not reachable\n", endpoint)
		}

		if err == nil {
			respOutput(resp, endpoint)
		}
	}
}
