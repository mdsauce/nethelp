package diagnostics

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// PublicSites attempts to prove that the machine has internet
// connectivity and is not being blocked by a private network.
func PublicSites(sitelist []string) {
	for _, site := range sitelist {
		log.Debug("Sending GET req to ", site)
		resp, err := http.Get(site)
		if err != nil {
			fmt.Printf("[ ] %s not reachable\n", site)
			log.WithFields(log.Fields{
				"error": err,
			}).Infof("[ ] %s not reachable\n", site)
		}
		if err == nil {
			respOutput(resp, site)
		}
	}
}
