package server

import (
	"fmt"
	"net/http"
)

const prefixLen = len("/players/")

// PlayerStore .
type PlayerStore interface {
	GetPlayerScore(name string) string
}

// PlayerServer .
type PlayerServer struct {
	Store PlayerStore
}

// StubPlayerStore 存储玩家分数
type StubPlayerStore struct {
	scores map[string]string
}

// GetPlayerScore .
func (s *StubPlayerStore) GetPlayerScore(name string) string {
	return s.scores[name]
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[prefixLen:]
	score := p.Store.GetPlayerScore(player)

	if score == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, p.Store.GetPlayerScore(player))
}
