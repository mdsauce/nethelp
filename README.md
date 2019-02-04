# nethelp
Nethelp will assist with finding out what is blocking outbound connections from the machine by sending HTTP and TCP connections to servies used by Sauce Labs.

### Usage
```
  nethelp [flags]

Flags:
      --config string   config file path (default is $HOME/.nethelp.yaml)
  -h, --help            help for nethelp
      --http            run HTTP tests. Default is to run all tests.
      --log             enables logging to the file specified by the --out flag.
  -l, --lucky           feeling lucky?  Disable the proxy check at startup and find out if it works during runtime.
  -o, --out string      optional output file for logging. Defaults to timestamp file in the current dir.  Only use if you want a custom log name. (default "20190203234732")
  -p, --proxy string    upstream proxy for nethelp to use. Enter like -p protocol://username:password@host:port
      --tcp             run TCP tests. Default is to run all tests.
```

### Build
Built using [Cobra](https://github.com/spf13/cobra) and go1.11.  Cobra is basically a templating tool for CLI and generator for the file structure. Cobra is built  on top of [pflag](https://github.com/spf13/pflag) which expands on the std library flag package in Go.

1. Clone the repo.
```
$ git clone git@github.com:mdsauce/nethelp.git
```
2. Go the nethelp dir.  Use `go build` to build a local version in the current dir or `go install` to install one in the go/bin folder and add the binary to your path.

If you're new to Go consider taking the tour https://tour.golang.org/list. 
