package server

import (
	"fmt"
	"net/http"
)

const prefixLen = len("/players/")

func PlayerServer(w http.ResponseWriter, r *http.Request) {

	player := r.URL.Path[prefixLen:]
	_, _ = fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(p string) string {
	if p == "Pepper" {
		return "10"
	}
	if p == "Floyd" {
		return "20"
	}
	return ""
}
