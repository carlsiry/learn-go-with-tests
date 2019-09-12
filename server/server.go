package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const prefixLen = len("/players/")

type Player struct {
	Name string
	Wins int
}

// PlayerStore 玩家数据存储器
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

// StubPlayerStore 玩家分数存储器
type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

// GetPlayerScore 获取玩家分数方法
func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

// InMemoryPlayerStore 内存实现版
type InMemoryPlayerStore struct {
	store map[string]int
}

// GetPlayerScore 获取玩家的分数
func (s *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return s.store[name]
}
func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.store[name]++
}
func (s *InMemoryPlayerStore) GetLeague() []Player {
	var league []Player
	for name, wins := range s.store {
		league = append(league, Player{
			Name: name,
			Wins: wins,
		})
	}
	return league
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{store: map[string]int{}}
}

// PlayerServer 玩家服务
type PlayerServer struct {
	Store PlayerStore
	http.Handler
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
