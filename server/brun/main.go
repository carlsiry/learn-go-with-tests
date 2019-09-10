package main

import (
	"net/http"

	"github.com/caicloud/nirvana/log"
	"github.com/carlsiry/learn-go-with-tests/server"
)

func main() {

	ps := &server.PlayerServer{
		Store: server.NewInMemoryPlayerStore(),
	}

	if err := http.ListenAndServe(":5000", ps); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
