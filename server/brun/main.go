package main

import (
	"net/http"

	"github.com/caicloud/nirvana/log"
	"github.com/carlsiry/learn-go-with-tests/server"
)

// InMemoryPlayerStore 内存实现版
type InMemoryPlayerStore struct{}

// GetPlayerScore 获取玩家的分数
func (s *InMemoryPlayerStore) GetPlayerScore(name string) string {
	return ""
}

func main() {

	server := &server.PlayerServer{&InMemoryPlayerStore{}}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
