package server

import (
	"bytes"
	"strings"
	"testing"
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

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		got := stdout.String()
		want := PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if game.StartedWith != 7 {
			t.Errorf("wanted start called with 7 but got %d", game.StartedWith)
		}
	})

	t.Run("it prints an error whenn a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessageSentToUser(t, stdout, PlayerPrompt, BadPlayerInputErrMsg)
		assertGameNotStarted(t, game)
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

func assertMessageSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()

	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func assertGameNotStarted(t testing.T, game *GameSpy) {
	t.Helper()
	if game.StartedWith != 0 {
		t.Errorf("game should not have started")
	}
}
