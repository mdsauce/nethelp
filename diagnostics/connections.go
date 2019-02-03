package diagnostics

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

//PublicSites attempts to send HTTP requests to sites that SHOULD be reachable.
func PublicSites(sitelist []string) {
	for _, site := range sitelist {
		log.Debug("Sending GET req to ", site)
		resp, err := http.Get(site)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatalf("[ ] %s not reachable\n", site)
		}

		if resp.StatusCode == 200 {
			fmt.Printf("[\u2713] %s is reachable %s\n", site, resp.Status)
		}
		log.WithFields(log.Fields{
			"status": resp.Status,
			"resp":   resp,
		}).Infof("[\u2713] %s reachable.\n", site)
	}
}

// SauceServices sends HTTP requests to Sauce endpoints
func SauceServices(sauceEndpoints []string) {
	for _, endpoint := range sauceEndpoints {
		log.Debug("Sending GET req to ", endpoint)
		resp, err := http.Get(endpoint)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatalf("[ ] %s not reachable\n", endpoint)
		}

		if resp.StatusCode == 200 {
			fmt.Printf("[\u2713] %s is reachable %s\n", endpoint, resp.Status)
			log.WithFields(log.Fields{
				"status": resp.Status,
				"resp":   resp,
			}).Infof("[\u2713] %s reachable.\n", endpoint)
		} else {
			fmt.Printf("[ ] %s returned %s\n", endpoint, resp.Status)
			log.WithFields(log.Fields{
				"status": resp.Status,
				"resp":   resp,
			}).Infof("[ ] %s returned %s\n", endpoint, resp.Status)
		}
	}
}

// RDCServices makes connections to the main RDC endpoints required to run tests
func RDCServices(rdcEndpoints []string) {
	for _, endpoint := range rdcEndpoints {
		log.Debug("Sending POST req to ", endpoint)
		var jsonBody = []byte(`{"test":"this will result in an HTTP 500 resp or 401 resp."}`)
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatalf("[ ] %s not reachable\n", endpoint)
		}

		if resp.StatusCode == 200 || resp.StatusCode == 401 || resp.StatusCode == 500 {
			fmt.Printf("[\u2713] %s is reachable %s\n", endpoint, resp.Status)
			log.WithFields(log.Fields{
				"status": resp.Status,
				"resp":   resp,
			}).Infof("[\u2713] %s reachable.\n", endpoint)
		} else {
			fmt.Printf("[ ] %s returned %s\n", endpoint, resp.Status)
			log.WithFields(log.Fields{
				"status": resp.Status,
				"resp":   resp,
			}).Infof("[ ] %s returned %s\n", endpoint, resp.Status)
		}
	}
}

//TCPConns attempts to open various TCP connections to the provided sites
func TCPConns(sitelist []string, proxyURL *url.URL) {
	if proxyURL != nil {
		log.Warn("May fail if you are not using a SOCKS5 Proxy.")
		var err error
		var proxyDialer proxy.Dialer
		proxyDialer, err = proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			log.Fatalf("Something went wrong while starting a proxy dialer for TCP conns.\n%v", err)
		}
		for _, site := range sitelist {
			conn, err := proxyDialer.Dial("tcp4", site)
			if err != nil {
				log.Fatalf("%s unreachable, %v: ", site, err)
			}
			fmt.Println("[\u2713] TCP (IPv4) connection to", site)
			log.WithFields(log.Fields{
				"local":  conn.LocalAddr(),
				"remote": conn.RemoteAddr(),
			}).Infof("[\u2713] %s reachable via TCP (IPv4).\n", site)
			conn.Close()
		}
	} else {
		for _, site := range sitelist {
			timeout := time.Duration(5 * time.Second)
			conn, err := net.DialTimeout("tcp4", site, timeout)
			if err != nil {
				log.Fatalf("%s unreachable, %v: ", site, err)
			}
			fmt.Println("[\u2713] TCP (IPv4) connection to", site)
			log.WithFields(log.Fields{
				"local":  conn.LocalAddr(),
				"remote": conn.RemoteAddr(),
			}).Infof("[\u2713] %s reachable via TCP (IPv4).\n", site)
			conn.Close()
		}
	}
}
