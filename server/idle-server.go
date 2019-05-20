package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i, err := strconv.ParseInt(r.URL.Path[1:], 10, 32)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("received request with a idle of %d seconds\n", i)

		time.Sleep(time.Duration(i) * time.Second)
		fmt.Printf("Waited for %d seconds\n", i)
		fmt.Printf("request answered after %d seconds\n", i)
	})
}

// IdleServer starts a webserver that will wait 15+ minutes
// before responding to HTTP requests
func IdleServer() {
	port := "8080"

	s := &http.Server{
		Addr:         fmt.Sprint(":", port),
		Handler:      handler(),
		WriteTimeout: 20 * time.Minute,
		IdleTimeout:  20 * time.Minute,
	}

	fmt.Println("Starting server, listening on port ", port)
	fmt.Println(s.ListenAndServe())
}
