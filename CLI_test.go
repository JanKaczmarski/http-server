package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/jankaczmarski/http-server"
)

func TestCLI(t *testing.T) {
	var dummyStdOut = &bytes.Buffer{}

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &poker.GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessageSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishedCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &poker.GameSpy{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishedCalledWith(t, game, "Cleo")
	})

	t.Run("it prints error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &poker.GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("Pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
		assertGameNotStarted(t, game)
	})

	t.Run("it prints error when wrong message is sent instead of 'Player wins'", func(t *testing.T) {
		game := &poker.GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("5", "Lloyd is a killer")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputErrMsg)
		assertGameNotFinished(t, game)
	})
}

func assertGameNotFinished(t testing.TB, game *poker.GameSpy) {
	t.Helper()

	if game.FinishCalled {
		t.Errorf("game should not have finished")
	}

}

func assertFinishedCalledWith(t testing.TB, game *poker.GameSpy, wantWinner string) {
	t.Helper()

	if game.FinishedWith != wantWinner {
		t.Errorf("wanted Finished called with %s but got %s", game.FinishedWith, wantWinner)
	}
}

func assertGameStartedWith(t testing.TB, game *poker.GameSpy, wantStartNumberOfPlayers int) {
	t.Helper()

	if game.StartedWith != wantStartNumberOfPlayers {
		t.Errorf("wanted Start called with %d but got %d", wantStartNumberOfPlayers, game.StartedWith)
	}
}

func assertMessageSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func assertGameNotStarted(t testing.TB, game *poker.GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}
