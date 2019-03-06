package main

import (
    "flag"
    "fmt"
    "net/http"
    "time"
    "strconv"

	log "github.com/sirupsen/logrus"
)

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i, err := strconv.ParseInt(r.URL.Path[1:], 10, 32)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Debug("Received request for a timeout of ", i, " seconds")

		time.Sleep(time.Duration(i) * time.Second)
		fmt.Fprintln(w, "Waited for", i, "seconds")
		log.Debug("request answered after ", i, " seconds")
	})
}

func main() {
    port := flag.String("p", "8080", "port to listen to")
    verbose := flag.Bool("v", false, "verbose mode")
    flag.Parse()

	log.SetLevel(log.WarnLevel)

	s := &http.Server{
		Addr:           fmt.Sprint(":", *port),
		Handler:        handler(),
		WriteTimeout:   20 * time.Minute,
		IdleTimeout: 20 * time.Minute,
	}

	if *verbose {
		log.SetLevel(log.TraceLevel)
	}

	log.Info("Starting server, listening on port ", *port)
	log.Fatal(s.ListenAndServe())
}