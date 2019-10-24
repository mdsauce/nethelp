package endpoints

import (
	"fmt"
	"os"
)

// AssembleHeadlessEndpoints interpolates user variables like
// SAUCE_USERNAME and SAUCE_ACCESS_KEY to create a valid URI.
func AssembleHeadlessEndpoints(dc string) *SauceService {
	if os.Getenv("SAUCE_USERNAME") == "" {
		log.Info("SAUCE_USERNAME environment variables not found.  Not running VDC REST endpoint tests.")
		return nil
	}
	eastHeadless := fmt.Sprintf("https://us-east-1.saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))

	switch dc {
	case "all":
		e := make([]string, 1)
		e[0] = eastHeadless
		return &SauceService{Datacenter: dc, Cloud: "headless", Endpoints: e}
	case "east":
		e := make([]string, 1)
		e[0] = eastHeadless
		return &SauceService{Datacenter: dc, Cloud: "headless", Endpoints: e}
	default:
		return nil
	}
}

// NewHeadlessTest constructs a SauceService object that contains the specificed Datacenter and endpoints
func NewHeadlessTest(dc string) SauceService {
	headlessTest := SauceService{Datacenter: dc, Cloud: "headless"}
	if dc == "east" || dc == "all" {
		headlessTest.Endpoints = []string{"http://ondemand.us-east-1.saucelabs.com:80", "https://ondemand.us-east-1.saucelabs.com:443"}
	}
	return headlessTest
}
