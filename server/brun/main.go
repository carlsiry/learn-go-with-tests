package main

import (
	"net/http"

	"github.com/caicloud/nirvana/log"
	"github.com/carlsiry/learn-go-with-tests/server"
)

func main() {

	handler := http.HandlerFunc(server.PlayerServer)

	if err := http.ListenAndServe(":5000", handler); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
