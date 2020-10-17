package main

import "testing"

func TestHello(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()

		if got != want {
			t.Errorf("got '%q' want '%q'", got, want)
		}
	}

	t.Run("saying Hello to Faker", func(t *testing.T) {
		got := Hello("Faker")
		want := "Hello, Faker"
		assertCorrectMessage(t, got, want)
	})

	t.Run("empty string defaults  to 'world'", func(t *testing.T) {
		got := Hello("")
		want := "Hello, world"
		assertCorrectMessage(t, got, want)
	})
}
