package main

import (
	"sync"
	"testing"
)

func TestInMemoeryPlayerStoreGetPlayerScore(t *testing.T) {
	store := NewInMemoryPlayerStore()

	t.Run("get active player score", func(t *testing.T) {
		player := "Bob"
		store.store[player] = 1

		got := store.GetPlayerScore(player)
		want := 1

		assertScore(t, got, want)
	})
}

func TestInMemoryPlayerStoreRecordWin(t *testing.T) {

	t.Run("record single win", func(t *testing.T) {
		store := NewInMemoryPlayerStore()
		player := "Bob"

		store.RecordWin(player)

		got := store.GetPlayerScore(player)
		want := 1

		assertScore(t, got, want)
	})

	t.Run("record many wins concurenntly", func(t *testing.T) {
		store := NewInMemoryPlayerStore()
		player := "Bob"
		recordCount := 1000

		var wg sync.WaitGroup
		wg.Add(recordCount)

		for i := 0; i < recordCount; i++ {
			go func() {
				store.RecordWin(player)

				wg.Done()
			}()
		}

		wg.Wait()

		got := store.GetPlayerScore(player)

		assertScore(t, got, recordCount)
	})
}

func assertScore(t testing.TB, got, want int) {
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
