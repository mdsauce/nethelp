[![CircleCI](https://circleci.com/gh/mdsauce/nethelp.svg?style=svg)](https://circleci.com/gh/mdsauce/nethelp)

# nethelp
Nethelp will assist with finding out what is blocking outbound connections from the machine by sending HTTP and TCP connections to servies used by Sauce Labs.

### Downloading and setting u
1. Download the binary for your operating system from [the releases](https://github.com/mdsauce/nethelp/releases).  
2. On Mac and Linux make this file executable by running `$ chmod 755`. You may get a `permission denied` type error if you try and run without the `chmod` step.   For example on a Linux machine:
```
$ cd ~/Downloads/nethelp-linux
$ chmod 755 nethelp
$ ./nethelp --help
```
3. Run `./nethelp` with whatever flags you need.
[![asciicast](https://asciinema.org/a/232600.svg)](https://asciinema.org/a/232600)

If you are on a Linux or Mac OS and you have root access you can add the `nethelp` binary to your command line by moving it to `/usr/local/bin`.  

You can also export the file `export PATH=$PATH:</path/to/file>` by adding that line to your `~/.bashrc` or `~/.bash_profile`.  More information here: https://unix.stackexchange.com/questions/3809/how-can-i-make-a-program-executable-from-everywhere.

### Usage
```
$ nethelp --help

 ___  __ _ _   _  ___ ___   / / __   ___| |_| |__   ___| |_ __  
/ __|/ _  | | | |/ __/ _ \ / / '_ \ / _ \ __| '_ \ / _ \ | '_ \ 
\__ \ (_| | |_| | (_|  __// /| | | |  __/ |_| | | |  __/ | |_) |
|___/\__,_|\__,_|\___\___/_/ |_| |_|\___|\__|_| |_|\___|_| .__/ 
                                                         |_|  
Nethelp will help find out what is blocking outbound 
connections by sending requests to 
services used during typical Sauce Labs usage.

Usage:
  nethelp [flags]

Flags:
      --api            run API tests.  Requires that you have $SAUCE_USERNAME and $SAUCE_ACCESS_KEY environment variables.
      --cloud string   options are: VDC or RDC.  Select which services you'd like to test, Virtual Device Cloud or Real Device Cloud respectively. (default "all")
      --dc string      options are: EU or NA.  Choose which data centers you want run diagnostics against, Europe or North America respectively. (default "all")
  -h, --help           help for nethelp
      --log            enables logging and creates a nethelp.log file.  Will automatically append data to the file in a non-destructive manner.
  -l, --lucky          disable the proxy check at startup and instead test the proxy during execution.
  -p, --proxy string   upstream proxy for nethelp to use. Enter like -p protocol://username:password@host:port
      --tcp            run TCP tests. Will always run against all endpoints.
  -v, --verbose        print all logging levels
```

* Run HTTP and API tests with a proxy upstream from your machine
```
$ nethelp --api --http -p myUsername:myPassword@upstream.proxy.inc.com:8080

```

* Log in Verbose mode and save to a logfile
```
$ nethelp -v --log
```

* Disable the initial proxy validation
```
$ nethelp -l
```

* Run tests only against a specific data center and cloud service
```
$ nethelp  --cloud vdc --dc na
[✓] https://ondemand.saucelabs.com:443 is reachable 200 OK
[✓] http://ondemand.saucelabs.com:80 is reachable 200 OK
```

### Idle server (for development only)

**This server is not needed client side, it is only necessary for the Sauce Labs support team**

The repository provides a custom http server simulating long idle connections. It can be found in `server/idle-server.go`.
Usage:

```
$ cd server
$ go build idle-server.go
$ ./idle-server -p 8080 -v
INFO[0000] Starting server, listening on port 8080
```

You can request specific timeouts in seconds the following way:

```
$ curl http://localhost:8080/10 # 10 seconds timeout
$ curl http://localhost:8080/900 # 15 minutes timeout
```

The server will answer after the requested number of seconds, allowing to simulate long running idle connections.
This is especially useful when trying to find out if long allocation time for RDC is a problem from a specific network.

### Build
Built using [Cobra](https://github.com/spf13/cobra) and go v1.11.  Cobra is an opinionated CLI generator. Cobra is built  on top of [pflag](https://github.com/spf13/pflag) which expands on the std library flag package in Go.

1. Clone the repo.
```
$ git clone git@github.com:mdsauce/nethelp.git
```
2. Go the nethelp dir.  Use `$ go build` to build a local version in the current dir or `$ go install` to install one in the `~/go/bin` folder and add the `nethelp` binary to your path.

If you're new to Go consider taking the tour https://tour.golang.org/list.

### Next Features
* create a test session with a specific name then quit it.  This will prove a connection can be made to the services and a test can start with user credentials.
* recover from failures automatically, record the failure, then continue the rest of the diagnostics.
