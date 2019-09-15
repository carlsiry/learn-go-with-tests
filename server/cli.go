package server

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const PlayerPrompt = "please enter the number of players: "

type CLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
	out         io.Writer
	alerter     BlindAlerter
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) readLine() string {
	line, _, _ := cli.in.ReadLine()
	return string(line)
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Second
	}
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
