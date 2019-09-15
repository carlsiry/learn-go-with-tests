package server

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type CLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
	out         io.Writer
	alerter     BlindAlerter
}

func (cli *CLI) PlayPoker() {
	cli.alerter.ScheduleAlertAt(5*time.Second, 100)
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) readLine() string {
	line, _, _ := cli.in.ReadLine()
	return string(line)
}

func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewReader(in),
		out:         out,
		alerter:     alerter,
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
