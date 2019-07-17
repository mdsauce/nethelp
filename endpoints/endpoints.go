package endpoints

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Check is the target of endpoints that
// should be reachable
type Check struct {
	Sitelist []string
}

// SauceService is a combination of the
// Cloud, Geographic location of the DC, and endpoint collection
type SauceService struct {
	Datacenter string
	Cloud      string
	Endpoints  []string
}

// vdcNA = []string{"https://ondemand.saucelabs.com:443", "http://ondemand.saucelabs.com:80"}
// vdcEU = []string{"http://ondemand.eu-central-1.saucelabs.com:80", "https://ondemand.eu-central-1.saucelabs.com:443"}
// rdcNA = []string{"https://us1.appium.testobject.com/wd/hub/session"}
// rdcEU = []string{"https://eu1.appium.testobject.com/wd/hub/session"}

// NewTCPTest builds a new TCPTest object
func NewTCPTest() Check {
	defaultTCP := Check{}
	defaultTCP.Sitelist = []string{"ondemand.saucelabs.com:443", "ondemand.saucelabs.com:80", "ondemand.saucelabs.com:8080", "ondemand.eu-central-1.saucelabs.com:80", "ondemand.eu-central-1.saucelabs.com:443", "us1.appium.testobject.com:443", "eu1.appium.testobject.com:443", "us1.appium.testobject.com:80", "eu1.appium.testobject.com:80"}
	return defaultTCP
}

// NewPublicTest assembles a collection of default endpoints
// that should all be reachable to assert connectivity to the public internet
func NewPublicTest() Check {
	defaultPublic := Check{}
	defaultPublic.Sitelist = []string{"https://status.us-west-1.saucelabs.com", "http://status.eu-central-1.saucelabs.com/", "https://www.duckduckgo.com"}
	return defaultPublic
}

// NewRDCTest takes a Data Center and assembles a collection of endpoints
// and geographic + service definitions
func NewRDCTest(dc string) SauceService {
	rdcTest := SauceService{Datacenter: dc, Cloud: "rdc"}
	if dc == "eu" {
		rdcTest.Endpoints = []string{"https://eu1.appium.testobject.com/wd/hub/session"}
	}
	if dc == "na" {
		rdcTest.Endpoints = []string{"https://us1.appium.testobject.com/wd/hub/session"}
	}
	if dc == "all" {
		rdcTest.Endpoints = []string{"https://eu1.appium.testobject.com/wd/hub/session", "https://us1.appium.testobject.com/wd/hub/session"}
	}
	return rdcTest
}

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

// NewHeadlessTest constructs a SauceService object that contains the specificed Datacenter and endpoints
func NewHeadlessTest(dc string) SauceService {
	headlessTest := SauceService{Datacenter: dc, Cloud: "headless"}
	if dc == "east" {
		headlessTest.Endpoints = []string{"http://ondemand.us-east-1.saucelabs.com:80", "https://ondemand.us-east-1.saucelabs.com:443"}
	}
	return headlessTest
}

// AssembleVDCEndpoints interpolates user variables like
// SAUCE_USERNAME and SAUCE_ACCESS_KEY to create a valid URI.
func AssembleVDCEndpoints(dc string) (*SauceService, error) {
	if os.Getenv("SAUCE_USERNAME") == "" {
		log.Info("SAUCE_USERNAME environment variables not found.  Not running VDC REST endpoint tests.")
		return nil, errors.New("SAUCE_USERNAME environment variables not found, not running VDC REST endpoint tests")
	}
	naEndpoint := fmt.Sprintf("https://saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))
	euEndpoint := fmt.Sprintf("https://eu-central-1.saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))

	switch dc {
	case "all":
		e := make([]string, 2)
		e[0] = naEndpoint
		e[1] = euEndpoint
		return &SauceService{Datacenter: dc, Cloud: "vdc", Endpoints: e}, nil
	case "na":
		e := make([]string, 1)
		e[0] = naEndpoint
		return &SauceService{Datacenter: dc, Cloud: "vdc", Endpoints: e}, nil
	case "eu":
		e := make([]string, 1)
		e[0] = euEndpoint
		return &SauceService{Datacenter: dc, Cloud: "vdc", Endpoints: e}, nil
	default:
		return nil, errors.New("Only 'all', 'vdc', or 'rdc' are allowed")
	}
}

// AssembleHeadlessEndpoints interpolates user variables like
// SAUCE_USERNAME and SAUCE_ACCESS_KEY to create a valid URI.
func AssembleHeadlessEndpoints(dc string) (*SauceService, error) {
	if os.Getenv("SAUCE_USERNAME") == "" {
		log.Info("SAUCE_USERNAME environment variables not found.  Not running VDC REST endpoint tests.")
		return nil, errors.New("SAUCE_USERNAME environment variables not found, not running VDC REST endpoint tests")
	}
	westHeadless := fmt.Sprintf("https://us-east-1.saucelabs.com/rest/v1/%s/tunnels", os.Getenv("SAUCE_USERNAME"))

	switch dc {
	case "all":
		e := make([]string, 1)
		e[0] = westHeadless
		return &SauceService{Datacenter: dc, Cloud: "headless", Endpoints: e}, nil
	case "east":
		e := make([]string, 1)
		e[0] = westHeadless
		return &SauceService{Datacenter: dc, Cloud: "headless", Endpoints: e}, nil
	default:
		return nil, errors.New("Only 'east' or 'all' is allowed")
	}
}
