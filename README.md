# nethelp
Nethelp will assist with finding out what is blocking outbound connections from the machine by sending HTTP and TCP connections to servies used by Sauce Labs.

### Help
```
Usage:
  nethelp [flags]

Flags:
      --api            run API tests.  Requires that you have $SAUCE_USERNAME and $SAUCE_ACCESS_KEY environment variables.
  -h, --help           help for nethelp
      --http           run HTTP tests. Default is to run all tests.
      --log            enables logging to the file specified by the --out flag.
  -l, --lucky          disable the proxy check at startup and instead test the proxy during execution.
  -p, --proxy string   upstream proxy for nethelp to use. Enter like -p protocol://username:password@host:port
      --tcp            run TCP tests. Default is to only run HTTP tests.
  -v, --verbose        print all logging levels
```

```
./nethelp-mac
[✓] https://status.saucelabs.com is reachable 200 OK
[✓] https://www.duckduckgo.com is reachable 200 OK
[✓] https://ondemand.saucelabs.com:443 is reachable 200 OK
[✓] http://ondemand.saucelabs.com:80 is reachable 200 OK
[✓] https://us1.appium.testobject.com/wd/hub/session is reachable but returned 401 Unauthorized
[✓] https://eu1.appium.testobject.com/wd/hub/session is reachable but returned 401 Unauthorized
```

### Downloading and using
Download the binary for your operating system at https://github.com/mdsauce/nethelp/releases.
On Mac and Linux make this file executable by running `$ chmod 755`.  For example:
```
$ chmod 755 nethelp-linux
```
You may get a `permission denied` type error if you try and run without this step.

### Build
Built using [Cobra](https://github.com/spf13/cobra) and go1.11.  Cobra is basically a templating tool for CLI and generator for the file structure. Cobra is built  on top of [pflag](https://github.com/spf13/pflag) which expands on the std library flag package in Go.

1. Clone the repo.
```
$ git clone git@github.com:mdsauce/nethelp.git
```
2. Go the nethelp dir.  Use `go build` to build a local version in the current dir or `go install` to install one in the go/bin folder and add the binary to your path.

If you're new to Go consider taking the tour https://tour.golang.org/list. 

### Next Features
* run VDC or RDC tests only with a `cloud` flag
* specify a data center, EU or US with a `dc` flag
* create a test session with a specific name then quit it.  This will prove a connection can be made to the services and a test can start.