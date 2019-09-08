package maps

import "testing"

func TestSearch(t *testing.T) {

	dictionary := Dictionary{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		got, _ := dictionary.Search("test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, err := dictionary.Search("unknown")
		want := ErrNotFound

		if err == nil {
			t.Fatalf("expected to get an error")
		}

		assertStrings(t, err.Error(), want.Error())
	})
}

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{}

		err := dictionary.Add(word, definition)

		assertError(t, err, nil)
		assertDefinition(t, dictionary, word, definition)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}

		err := dictionary.Add(word, definition)
		assertError(t, err, ErrWordExists)
		assertDefinition(t, dictionary, word, definition)
	})
}

func assertStrings(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error is '%s' want '%s'", got, want)
	}
}

func assertDefinition(t *testing.T, dict Dictionary, key, val string) {
	t.Helper()

	got, err := dict.Search(key)
	if err != nil {
		t.Fatal("should find added world, but", err)
	}
	if got != val {
		t.Errorf("got %s want %s", got, val)
	}
}
