package server

import (
	"encoding/json"
	"io"
)

type FileSystemPlayerStore struct {
	Database io.ReadWriteSeeker
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.GetLeague().Find(name)
	if player != nil {
		return player.Wins
	}
	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{
			Name: name,
			Wins: 1,
		})
	}

	f.Database.Seek(0, 0)
	json.NewEncoder(f.Database).Encode(league)
}

func (f *FileSystemPlayerStore) GetLeague() League {
	_, _ = f.Database.Seek(0, 0)
	league, _ := NewLeague(f.Database)
	return league
}
