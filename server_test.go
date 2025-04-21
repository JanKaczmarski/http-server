package poker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var (
	dummyGame = &GameSpy{}
)

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
		sync.Mutex{},
	}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing player", func(t *testing.T) {
		request := NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		if got != want {
			t.Errorf("got status %d want %d", got, want)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns 200 on /league", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague, sync.Mutex{}}
		server := mustMakePlayerServer(t, &store, dummyGame)

		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetLeagueRequest(t, response.Body)
		AssertStatus(t, response, http.StatusOK)
		AssertLeague(t, got, wantedLeague)
		AssertContentType(t, response, jsonContentType)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
		sync.Mutex{},
	}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		player := "Pepper"

		request := NewPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusAccepted)
		AssertPlayerWin(t, &store, player)

	})

	t.Run("it runs concurrently POST", func(t *testing.T) {
		player := "Bob"
		wantedCalls := 1000

		var wg sync.WaitGroup
		wg.Add(wantedCalls)

		for i := 0; i < wantedCalls; i++ {
			go func() {
				request := NewPostWinRequest(player)
				response := httptest.NewRecorder()

				server.ServeHTTP(response, request)

				wg.Done()
			}()
		}

		wg.Wait()

		if len(store.winCalls) < wantedCalls {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), wantedCalls)
		}

	})
}

func TestGame(t *testing.T) {
	dummyPlayerStore := &StubPlayerStore{}
	t.Run("Get /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusOK)
	})

	t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
		game := &GameSpy{}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)

		if game.StartedWith != 3 {
			t.Fatalf("game staretd with wrong number of players want 3 got %d", game.StartedWith)
		}

		if game.FinishedWith != winner {
			t.Fatalf("game finished with wrong winner want %s got %s", winner, game.FinishedWith)
		}

	})
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)

	return request
}

func mustMakePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatalf("problem creating player server %v", err)
	}

	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
