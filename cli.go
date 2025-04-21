package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
const BadWinnerInputErrMsg = "Bad value received for winning command, pleas try again with 'PLAYER_NAME wins'"
const PlayerPrompt = "Please enter the number of players: "

func (cli *CLI) PlayPoker() {
	fmt.Fprintf(cli.out, PlayerPrompt)

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	if err != nil {
		fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		return
	}

	cli.game.Start(numberOfPlayers, cli.out)

	winnerInput := cli.readLine()
	winner, err := extractWinner(winnerInput)
	if err != nil {
		fmt.Fprint(cli.out, BadWinnerInputErrMsg)
		return
	}

	cli.game.Finish(winner)
}

func extractWinner(userInput string) (string, error) {
	re := regexp.MustCompile(`^(\w+) wins\s*$`)
	if re.MatchString(userInput) {
		return strings.Replace(userInput, " wins", "", 1), nil
	}

	return "", errors.New("input did not match the expected pattern")
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
