package main

import (
	"fmt"
	"log"
	"net/http"
)

func startHTTPServer() *http.Server {
	srv := &http.Server{Addr: ":8181"}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//     io.WriteString(w, "hello world\n")
	// })

	http.HandleFunc("/", sayHello)

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf(`{
	"environments": {
		"ready": "yes"
	}
}`)
	w.Write([]byte(message))
}
