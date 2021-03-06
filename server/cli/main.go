package main

import (
	"fmt"
	"github.com/carlsiry/learn-go-with-tests/server"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, err := server.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatalf("problem createing file system player store, %v", err)
	}

	game := server.NewGame(server.BlindAlerterFunc(server.StdOutAlerter), store)

	cli := server.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()
}
