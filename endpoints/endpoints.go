package endpoints

// // Connector makes connection tests
// type Connector interface{}

//TCPTest is used to test raw TCP connections to a predefined sitelist
type TCPTest struct {
	Sitelist []string
}

// NewTCPTest builds a new TCPTest object
func NewTCPTest() TCPTest {
	defaultTCP := TCPTest{}
	defaultTCP.Sitelist = []string{"ondemand.saucelabs.com:443", "ondemand.saucelabs.com:80", "ondemand.saucelabs.com:8080", "ondemand.eu-central-1.saucelabs.com:80", "ondemand.eu-central-1.saucelabs.com:443", "us1.appium.testobject.com:443", "eu1.appium.testobject.com:443", "us1.appium.testobject.com:80", "eu1.appium.testobject.com:80"}
	return defaultTCP
}
