package racer

import (
	"fmt"
	"net/http"
	"time"
)

const timeout = 3 * time.Second

var TimeoutError = fmt.Errorf("http request timeout within %d millsecond", timeout)

func Racer(a, b string) (string, error) {
	return ConfigurableRacer(a, b, timeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, err error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", TimeoutError
	}
}

func ping(url string) chan bool {
	ch := make(chan bool)
	go func() {
		http.Get(url)
		ch <- true
	}()
	return ch
}
