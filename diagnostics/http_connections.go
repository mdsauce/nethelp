package diagnostics

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//PublicSites attempts to send HTTP requests to sites that SHOULD be reachable.
func PublicSites(sitelist []string) {
	for _, site := range sitelist {
		log.Debug("Sending GET req to ", site)
		resp, err := http.Get(site)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"resp":  resp,
			}).Errorf("[ ] %s not reachable %s\n", site, resp.Status)
		}

		if resp.StatusCode == 200 {
			fmt.Printf("[x] %s is reachable %s\n", site, resp.Status)
		}
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[x] %s reachable.\n", site)
	}

	// log.Debug("Sending GET req to www.duckduckgo.com.")
	// if resp, err = http.Get("https://www.duckduckgo.com"); err != nil {
	// 	log.WithFields(log.Fields{
	// 		"error": err,
	// 		"resp":  resp,
	// 	}).Error("[ ] www.duckduckgo.com not reachable", resp.Status)
	// }

	// if resp.StatusCode == 200 {
	// 	fmt.Println("[x] www.duckduckgo.com is reachable", resp.Status)
	// }
	// log.WithFields(log.Fields{
	// 	"status": resp.Status,
	// 	"resp":   resp,
	// }).Info("[x] www.duckduckgo.com reachable.")
}
