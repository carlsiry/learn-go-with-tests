package server

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

var (
	dummyPlayerStore  = &StubPlayerStore{}
	dummyBlindAlerter = &SpyBlindAlerter{}
	dummyStdIn        = strings.NewReader("5\n")
	dummyStdOut       = &bytes.Buffer{}
)

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("5\nChris wins\n")
		playerStore := &StubPlayerStore{}
		game := NewGame(dummyBlindAlerter, playerStore)

		cli := NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, dummyPlayerStore)

		cli := NewCLI(dummyStdIn, dummyStdOut, game)
		cli.PlayPoker()

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

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}

		game := NewGame(blindAlerter, dummyPlayerStore)
		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		got := stdout.String()
		want := PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Second, 200},
			{24 * time.Second, 300},
			{36 * time.Second, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})
}

func assertScheduledAlert(t *testing.T, got, want scheduledAlert) {
	t.Helper()
	amountGot := got.amount
	amountWant := want.amount
	if amountGot != amountWant {
		t.Errorf("got amount %d, want %d", amountGot, amountWant)
	}

	gotScheduledTime := got.at
	wantScheduledTime := want.at
	if gotScheduledTime != wantScheduledTime {
		t.Errorf("got amount %v, want %v", gotScheduledTime, wantScheduledTime)
	}
}

func assertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("didn't correct winner, got '%s', want '%s'", store.winCalls[0], winner)
	}
}
