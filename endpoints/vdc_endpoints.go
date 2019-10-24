package endpoints

import (
	"fmt"
	"os"
)

// NewVDCTest constructs a SauceService object that contains the specificed Datacenter and endpoints
func NewVDCTest(dc string) SauceService {
	vdcTest := SauceService{Datacenter: dc, Cloud: "vdc"}
	if dc == "eu" {
		vdcTest.Endpoints = []string{"http://ondemand.eu-central-1.saucelabs.com:80", "https://ondemand.eu-central-1.saucelabs.com:443"}
	}
	if dc == "na" {
		vdcTest.Endpoints = []string{"https://ondemand.saucelabs.com:443", "http://ondemand.saucelabs.com:80"}
	}
	if dc == "all" {
		vdcTest.Endpoints = []string{"http://ondemand.eu-central-1.saucelabs.com:80", "https://ondemand.eu-central-1.saucelabs.com:443", "https://ondemand.saucelabs.com:443", "http://ondemand.saucelabs.com:80"}
	}
	return vdcTest
}

// AssembleVDCEndpoints interpolates user variables like
// SAUCE_USERNAME and SAUCE_ACCESS_KEY to create a valid URI.
func AssembleVDCEndpoints(dc string) *SauceService {
	if os.Getenv("SAUCE_USERNAME") == "" {
		log.Warn("SAUCE_USERNAME environment variables not found.  Not running VDC REST endpoint tests.")
		return nil
	}
	naEndpoint := fmt.Sprintf("https://saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))
	euEndpoint := fmt.Sprintf("https://eu-central-1.saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))

	switch dc {
	case "all":
		e := make([]string, 2)
		e[0] = naEndpoint
		e[1] = euEndpoint
		return &SauceService{Datacenter: dc, Cloud: "vdc", Endpoints: e}
	case "na":
		e := make([]string, 1)
		e[0] = naEndpoint
		return &SauceService{Datacenter: dc, Cloud: "vdc", Endpoints: e}
	case "eu":
		e := make([]string, 1)
		e[0] = euEndpoint
		return &SauceService{Datacenter: dc, Cloud: "vdc", Endpoints: e}
	default:
		return nil
	}
}
