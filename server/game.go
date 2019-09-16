package server

import "time"

type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

type PokerGame struct {
	alerter BlindAlerter
	store   PlayerStore
}

func (p *PokerGame) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Second

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.ScheduleAlertAt(blindTime, blind)
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

type GameSpy struct {
	StartedWith  int
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
}

func (g GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}
