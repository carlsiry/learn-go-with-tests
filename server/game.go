package server

import (
	"io"
	"time"
)

type Game interface {
	Start(numberOfPlayers int, alertsDst io.Writer)
	Finish(winner string)
}

type PokerGame struct {
	alerter BlindAlerter
	store   PlayerStore
}

func (p *PokerGame) Start(numberOfPlayers int, alertsDst io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Second

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.ScheduleAlertAt(blindTime, blind, alertsDst)
		blindTime = blindTime + blindIncrement
	}
}

func (p *PokerGame) Finish(winner string) {
	p.store.RecordWin(winner)
}

func NewGame(alerter BlindAlerter, store PlayerStore) *PokerGame {

	return &PokerGame{
		alerter: alerter,
		store:   store,
	}
}
