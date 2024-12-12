package main

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	tt, err := myTime()
	if err != nil {
		t.Errorf("Oшибка")
		return
	}
	if tt.Second() != time.Now().Second() {
		t.Errorf("секунды не совпали")
	}
	if tt.Minute() != time.Now().Minute() {
		t.Errorf("минуты не совпали")
	}
	if tt.Hour() != time.Now().Hour() {
		t.Errorf("часы не совпали")
	}
}
