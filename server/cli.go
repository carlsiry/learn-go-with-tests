package server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const PlayerPrompt = "please enter the number of players: "

type CLI struct {
	in   *bufio.Reader
	out  io.Writer
	game *PokerGame
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, _ := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}

func (cli *CLI) readLine() string {
	line, _, _ := cli.in.ReadLine()
	return string(line)
}

func NewCLI(in io.Reader, out io.Writer, game *PokerGame) *CLI {
	return &CLI{
		in:   bufio.NewReader(in),
		out:  out,
		game: game,
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
