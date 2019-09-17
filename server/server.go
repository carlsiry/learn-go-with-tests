package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
)

const prefixLen = len("/players/")

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// PlayerServer 玩家服务
type PlayerServer struct {
	Store PlayerStore
	http.Handler
	template *template.Template
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

func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)
	_, winnerMsg, _ := conn.ReadMessage()
	p.Store.RecordWin(string(winnerMsg))
}

func NewPlayerServer(store PlayerStore, file string) (*PlayerServer, error) {
	server := new(PlayerServer)

	tmpl, err := template.ParseFiles(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading template %s", err.Error())
	}

	router := http.NewServeMux()
	router.Handle("/players/", http.HandlerFunc(server.playersHandle))
	router.Handle("/league", http.HandlerFunc(server.leagueHandle))
	router.Handle("/game", http.HandlerFunc(server.game))
	router.Handle("/ws", http.HandlerFunc(server.webSocket))

	server.Store = store
	server.template = tmpl
	server.Handler = router

	return server, nil
}
