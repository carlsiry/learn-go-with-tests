package main

import (
	"github.com/caicloud/nirvana/log"
	"github.com/carlsiry/learn-go-with-tests/server"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {

	store, err := server.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatalf("didnt expect an error but got one, %v", err)
	}

	ps, err := server.NewPlayerServer(store, "./server/game.html")
	if err != nil {
		log.Fatalf("didnt expect an error but got one, %v", err)
	}

	if err := http.ListenAndServe(":5000", ps); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
