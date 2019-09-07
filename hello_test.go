package main

import "testing"

func TestHello(t *testing.T) {

	assertCorrectMessage := func(t *testing.T, got, want string) {
		// t.Helper() 需要告诉测试套件这个方法是辅助函数（helper）
		// 当测试失败时所报告的行号将在函数调用中而不是在辅助函数内部
		// 这将帮助其他开发人员更容易地跟踪问题
		t.Helper()
		if got != want {
			t.Errorf("got '%s', want '%s' !", got, want)
		}
	}

	t.Run("saying hey to people", func(t *testing.T) {
		got := Hello("Carl", "")
		want := "Hey, Carl"
		assertCorrectMessage(t, got, want)
	})

	t.Run("say 'hey, world' when an empty string is supplied", func(t *testing.T) {
		got := Hello("", "")
		want := "Hey, world"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("Elodie", "French")
		want := "Bonjour, Elodie"
		assertCorrectMessage(t, got, want)
	})

}
