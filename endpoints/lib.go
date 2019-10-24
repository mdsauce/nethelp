package endpoints

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
