package server

// PlayerStore 玩家数据存储器
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

// StubPlayerStore 玩家分数存储器
type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

// GetPlayerScore 获取玩家分数方法
func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}
