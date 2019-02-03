package diagnostics

import (
	log "github.com/sirupsen/logrus"
)

//PublicSites attempts to send HTTP requests to sites that SHOULD be reachable.
func PublicSites() {
	log.Debug("Sending GET req to status.saucelabs.com.")
}
