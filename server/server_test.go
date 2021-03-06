package server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	dummyGame = &GameSpy{}
)

const jsonContentType = "application/json"
const testHTML = "./game.html"

func TestPlayerServer(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"Pepper": 10,
			"Floyd":  20,
		},
	}
	server := mustMakePlayerServer(t, store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetStoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})

	t.Run("return Floyd's score", func(t *testing.T) {
		request := newGetStoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})

	t.Run("returns 404 on missing player", func(t *testing.T) {
		request := newGetStoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertStatus(t, got, want)
	})
}

func newGetStoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/players/"+name, nil)
	return request
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestStoreWins(t *testing.T) {
	player := "Pepper"
	store := &StubPlayerStore{scores: map[string]int{}}
	server, err := NewPlayerServer(store, testHTML, dummyGame)
	assertError(t, err, nil)

	t.Run("it records wins when POST", func(t *testing.T) {
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		gotWinner := store.winCalls[0]
		if gotWinner != player {
			t.Errorf("didnot store correct winner got '%s' want '%s'", gotWinner, player)
		}
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "[]")
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)
	assertError(t, err, nil)

	server, err := NewPlayerServer(store, testHTML, dummyGame)
	assertError(t, err, nil)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetStoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})

}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as json", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Test", 14},
		}

		store := &StubPlayerStore{league: wantedLeague}

		server, err := NewPlayerServer(store, testHTML, dummyGame)
		assertError(t, err, nil)

		request := newLeagueRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertContentType(t, response, jsonContentType)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
	})
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}
func getLeagueFromResponse(t *testing.T, body io.Reader) []Player {
	t.Helper()
	// todo: handle err.
	league, _ := NewLeague(body)
	return league
}

func assertLeague(t *testing.T, got, want League) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Header().Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.HeaderMap)
	}
}

func mustMakePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
	server, err := NewPlayerServer(store, testHTML, game)
	if err != nil {
		t.Fatalf("problem creating player server: %v", err)
	}
	return server
}

func TestFileSystemStore(t *testing.T) {
	t.Run("/league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertError(t, err, nil)

		got := store.GetLeague()
		want := League{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}
		assertLeague(t, got, want)
	})

	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertError(t, err, nil)

		got := store.GetLeague()
		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("/get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertError(t, err, nil)

		got := store.GetPlayerScore("Chris")
		want := 33

		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertError(t, err, nil)

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for not existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertError(t, err, nil)

		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		assertScoreEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		assertError(t, err, nil)
	})
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()
	tmpfile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestGame(t *testing.T) {

	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		assertFinishCalledWith(t, game, winner)
		assertGameStartedWith(t, game, 3)

		within(t, 10*time.Millisecond, func() {
			assertWebsocketGotMsg(t, ws, wantedBlindAlert)
		})
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		game := &GameSpy{}

		out := &bytes.Buffer{}
		in := userSends("3", "Chris wins")

		NewCLI(in, out, game).PlayPoker()

		assertMessageSentToUser(t, out, PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}
	return ws
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send a message over ws connection %v", err)
	}
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)
	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Errorf("time out")
	case <-done:

	}
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf("got %s, want %s", string(msg), want)
	}
}
