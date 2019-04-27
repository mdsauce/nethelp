//+build !windows

package diagnostics

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func respOutput(resp *http.Response, endpoint string) {
	if resp.StatusCode == 200 {
		fmt.Printf("[\u2713] %s is reachable %s\n", endpoint, resp.Status)
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[\u2713] %s reachable.\n", endpoint)
	} else if resp.StatusCode == 401 {
		fmt.Printf("[\u2713] %s is reachable but returned %s\n", endpoint, resp.Status)
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[\u2713] %s reachable but unauthenticated.\n", endpoint)
	} else {
		fmt.Printf("[ ] %s returned %s\n", endpoint, resp.Status)
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[ ] %s returned %s\n", endpoint, resp.Status)
	}
}
