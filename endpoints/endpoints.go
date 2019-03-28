package endpoints

// Check is used to assert that the listed endpoints are reachable
type Check struct {
	Sitelist []string
}

// Service is used to assert that the specified datacenter and cloud is reachable
type Service struct {
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

// NewPublicTest builds a new publicTest object for connectivity assertions to public websites
func NewPublicTest() Check {
	defaultPublic := Check{}
	defaultPublic.Sitelist = []string{"https://status.us-west-1.saucelabs.com", "http://status.eu-central-1.saucelabs.com/", "https://www.duckduckgo.com"}
	return defaultPublic
}

// NewRDCTest constructs a Service object that contains the specificed Datacenter and endpoints
func NewRDCTest(dc string) Service {
	rdcTest := Service{Datacenter: dc, Cloud: "rdc"}
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

// NewVDCTest constructs a Service object that contains the specificed Datacenter and endpoints
func NewVDCTest(dc string) Service {
	vdcTest := Service{Datacenter: dc, Cloud: "vdc"}
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
