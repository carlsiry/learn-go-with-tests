package main

import (
	"net/http"
	"os"

	"github.com/caicloud/nirvana/log"
	"github.com/carlsiry/learn-go-with-tests/server"
)

const dbFileName = "game.db.json"

func main() {

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store := &server.FileSystemPlayerStore{db}
	ps := server.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", ps); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
