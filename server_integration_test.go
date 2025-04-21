package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := CreateTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	AssertNoError(t, err)

	player := "Pepper"
	server := mustMakePlayerServer(t, store, dummyGame)

	if err != nil {
		t.Fatalf("Error creating new player server %v", err)
	}

	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewGetScoreRequest(player))
		AssertStatus(t, response, http.StatusOK)

		AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewLeagueRequest())
		AssertStatus(t, response, http.StatusOK)

		got := GetLeagueRequest(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		AssertLeague(t, got, want)
	})
}
