//+build windows

package connections

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func respOutput(resp *http.Response, endpoint string) {
	if resp.StatusCode == 200 {
		fmt.Printf("[OK] %s is reachable %s\n", endpoint, resp.Status)
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[OK] %s reachable.\n", endpoint)
	} else if resp.StatusCode == 401 {
		fmt.Printf("[OK] %s is reachable but returned %s\n", endpoint, resp.Status)
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[OK] %s reachable but unauthenticated.\n", endpoint)
	} else {
		fmt.Printf("[ERROR] %s returned %s\n", endpoint, resp.Status)
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[ERROR] %s returned %s\n", endpoint, resp.Status)
	}
}
