package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Carl")
	want := "Hey, Carl"

	if got != want {
		t.Errorf("got '%s', want '%s' !", got, want)
	}
}
