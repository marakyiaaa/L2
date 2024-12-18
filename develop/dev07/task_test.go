package main

import (
	"testing"
	"time"
)

func sig(after time.Duration) <-chan any {
	c := make(chan any)
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

// Тесты
func TestOr(t *testing.T) {
	t.Run("No channels", func(t *testing.T) {
		result := or()
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}
	})
	t.Run("Single channel", func(t *testing.T) {
		start := time.Now()
		<-or(
			sig(1 * time.Second),
		)
		elapsed := time.Since(start)
		if elapsed < 1*time.Second || elapsed > 1500*time.Millisecond {
			t.Errorf("Expected to close after 1 second, but took %v", elapsed)
		}
	})
}
