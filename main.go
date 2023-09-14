package main

import (
	"errors"
	"fmt"
	"net/http"
	"stddev-service/api"
)

func main() {
	runServer()
}

func runServer() {
	http.HandleFunc("/random/mean", api.GetRandomMean)
	err := http.ListenAndServe(":80", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
}
