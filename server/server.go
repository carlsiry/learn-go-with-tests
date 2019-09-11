package server

import (
	"fmt"
	"net/http"
)

const prefixLen = len("/players/")

// PlayerStore 玩家数据存储器
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

// StubPlayerStore 玩家分数存储器
type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

// GetPlayerScore 获取玩家分数方法
func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
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

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{store: map[string]int{}}
}

// PlayerServer 玩家服务
type PlayerServer struct {
	Store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()

	router.Handle("/players/", http.HandlerFunc(p.playersHandle))
	router.Handle("/league", http.HandlerFunc(p.leagueHandle))

	router.ServeHTTP(w, r)
}

func (p *PlayerServer) leagueHandle(w http.ResponseWriter, r *http.Request) {
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
