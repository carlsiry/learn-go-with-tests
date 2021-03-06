package server

import (
	"encoding/json"
	"fmt"
	"io"
)

type Player struct {
	Name string
	Wins int
}

type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

func NewLeague(rdr io.Reader) (league []Player, err error) {
	if err = json.NewDecoder(rdr).Decode(&league); err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}
	return
}
