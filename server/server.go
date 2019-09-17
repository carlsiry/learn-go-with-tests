package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const prefixLen = len("/players/")

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type playerServerWS struct {
	*websocket.Conn
}

func (w *playerServerWS) Write(p []byte) (n int, err error) {
	if err = w.WriteMessage(1, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("proble upgrading connection to websockets %v\n", err)
	}
	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}

// PlayerServer 玩家服务
type PlayerServer struct {
	Store PlayerStore
	http.Handler
	template *template.Template
	game     Game
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

func (p *PlayerServer) games(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	p.game.Start(numberOfPlayers, ws)

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}

func NewPlayerServer(store PlayerStore, file string, game Game) (*PlayerServer, error) {
	server := new(PlayerServer)

	tmpl, err := template.ParseFiles(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading template %s", err.Error())
	}

	router := http.NewServeMux()
	router.Handle("/players/", http.HandlerFunc(server.playersHandle))
	router.Handle("/league", http.HandlerFunc(server.leagueHandle))
	router.Handle("/game", http.HandlerFunc(server.games))
	router.Handle("/ws", http.HandlerFunc(server.webSocket))

	server.Store = store
	server.game = game
	server.template = tmpl
	server.Handler = router

	return server, nil
}
