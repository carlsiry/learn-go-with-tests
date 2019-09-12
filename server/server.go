package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const prefixLen = len("/players/")

// PlayerServer 玩家服务
type PlayerServer struct {
	Store PlayerStore
	http.Handler
}

func (p *PlayerServer) leagueHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// todo: handle error
	_ = json.NewEncoder(w).Encode(p.Store.GetLeague())

	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandle(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[prefixLen:]
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.Store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	_, _ = fmt.Fprint(w, p.Store.GetPlayerScore(player))
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)
	server.Store = store

	router := http.NewServeMux()
	router.Handle("/players/", http.HandlerFunc(server.playersHandle))
	router.Handle("/league", http.HandlerFunc(server.leagueHandle))

	server.Handler = router

	return server
}
