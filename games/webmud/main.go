package main

import (
	"net/http"
)

func main() {
	s := NewServer()

	hs := &http.Server{
		ReadHeaderTimeout: gTimeout,
		WriteTimeout:      gTimeout,
		Addr:              "127.0.0.1:8080",
		Handler:           s,
	}

	Log("Listening on %s", hs.Addr)
	if err := hs.ListenAndServe(); err != nil {
		panic(err)
	}
}
