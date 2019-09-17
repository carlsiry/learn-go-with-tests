package server

import (
	"fmt"
	"testing"
	"time"
)

func TestGame_Start(t *testing.T) {

	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5, nil)

		cases := []scheduledAlert{
			{at: 0 * time.Second, amount: 100},
			{at: 10 * time.Second, amount: 200},
			{at: 20 * time.Second, amount: 300},
			{at: 30 * time.Second, amount: 400},
			{at: 40 * time.Second, amount: 500},
			{at: 50 * time.Second, amount: 600},
			{at: 60 * time.Second, amount: 800},
			{at: 70 * time.Second, amount: 1000},
			{at: 80 * time.Second, amount: 2000},
			{at: 90 * time.Second, amount: 4000},
			{at: 100 * time.Second, amount: 8000},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7, nil)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Second, 200},
			{24 * time.Second, 300},
			{36 * time.Second, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})
}

func TestGame_Finish(t *testing.T) {
	store := &StubPlayerStore{}
	game := NewGame(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	assertPlayerWin(t, store, winner)
}

func checkSchedulingCases(t *testing.T, cases []scheduledAlert, blindAlerter *SpyBlindAlerter) {
	t.Helper()
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
}
