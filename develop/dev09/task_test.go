package main

import (
	"net/url"
	"os"
	"testing"
)

func TestNewWget(t *testing.T) {
	wget, err := NewWget()
	if err != nil {
		t.Fatalf("Ошибка при создании Wget: %v", err)
	}
	if wget == nil {
		t.Fatal("NewWget вернул nil")
	}
}

func TestParseURL(t *testing.T) {
	input := "https://example.com"
	parsed, err := url.Parse(input)
	if err != nil {
		t.Fatalf("Ошибка при разборе URL: %v", err)
	}
	if parsed.String() != input {
		t.Errorf("Ожидалось %s, получено %s", input, parsed.String())
	}
}

func TestCreateDir(t *testing.T) {
	dir := "./test_dir"
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		t.Fatalf("Ошибка при создании директории: %v", err)
	}
	err = os.Remove(dir)
	if err != nil {
		t.Fatalf("Ошибка при удалении директории: %v", err)
	}
}

func TestDownloadPage(t *testing.T) {
	wget, err := NewWget()
	if err != nil {
		t.Fatalf("Ошибка при создании Wget: %v", err)
	}

	wget.OutputDir = "./test_download"
	testURL := "https://example.com"

	err = wget.DownloadPage(testURL, 0)
	if err != nil {
		t.Errorf("Ошибка при скачивании страницы: %v", err)
	}
}
