package main

import (
	"os"
	"time"

	"github.com/carlsiry/learn-go-with-tests/mock"
)

type ConfigurableSleeper struct {
	duration time.Duration
}

func (c *ConfigurableSleeper) Sleep() {
	time.Sleep(c.duration)
}

func main() {
	mock.Countdown(os.Stdout, &ConfigurableSleeper{1 * time.Second})
}
