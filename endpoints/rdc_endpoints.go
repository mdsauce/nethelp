package endpoints

// NewRDCTest takes a Data Center and assembles a collection of endpoints
// and geographic + service definitions
func NewRDCTest(dc string) SauceService {
	rdcTest := SauceService{Datacenter: dc, Cloud: "rdc"}
	if dc == "eu" {
		rdcTest.Endpoints = []string{"https://eu1.appium.testobject.com/wd/hub/status"}
	}
	if dc == "na" {
		rdcTest.Endpoints = []string{"https://us1.appium.testobject.com/wd/hub/status"}
	}
	if dc == "all" {
		rdcTest.Endpoints = []string{"https://eu1.appium.testobject.com/wd/hub/status", "https://us1.appium.testobject.com/wd/hub/status"}
	}
	return rdcTest
}
