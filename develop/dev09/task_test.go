package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

func TestSaveFile(t *testing.T) {
	wget, _ := NewWget()
	path := "./test_output/testfile.txt"
	content := strings.NewReader("Test content")

	// Создаем каталог, если он не существует
	if err := os.MkdirAll("./test_output", 0755); err != nil {
		t.Fatalf("Ошибка при создании каталога: %v", err)
	}

	err := wget.saveFile(path, content)
	if err != nil {
		t.Fatalf("Ошибка при сохранении файла: %v", err)
	}

	defer os.RemoveAll("./test_output")

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	if string(data) != "Test content" {
		t.Errorf("Ожидался 'Test content', получено '%s'", string(data))
	}
}

func TestDownloadPageOne(t *testing.T) {
	wget, err := NewWget()
	if err != nil {
		t.Fatalf("Ошибка при создании Wget: %v", err)
	}
	wget.Depth = 1
	wget.OutputDir = "./test_download"
	testURL := "https://example.com"

	err = wget.DownloadPage(testURL, 0)
	if err != nil {
		t.Errorf("Ошибка при скачивании страницы: %v", err)
	}
}

// мок
func TestDownloadPageTwo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body><h1>Test Page</h1></body></html>"))
	}))
	defer ts.Close()

	parsedURL, _ := url.Parse(ts.URL)
	wget := &Wget{
		visit:     make(map[string]bool),
		OutputDir: "./test_output",
		URL:       parsedURL,
	}

	err := wget.DownloadPage(ts.URL, 0)
	if err != nil {
		t.Fatalf("Ошибка при загрузке страницы: %v", err)
	}

	path, _ := wget.createPath(ts.URL)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Файл %s не был создан", path)
	}
}

func TestCreatePath(t *testing.T) {
	wget, _ := NewWget()
	wget.OutputDir = "./test_output"

	path, err := wget.createPath("http://example.com/test")
	if err != nil {
		t.Fatalf("Ошибка создания пути: %v", err)
	}

	expectedPath := filepath.Join(wget.OutputDir, "example.com", "test.html")
	if path != expectedPath {
		t.Errorf("Ожидался путь %s, получен %s", expectedPath, path)
	}
}

func TestResolveURL(t *testing.T) {
	wget, _ := NewWget()

	baseURL := "http://example.com"
	link := "/test/page.html"

	expected := "http://example.com/test/page.html"
	result := wget.resolveURL(link, baseURL)

	if result != expected {
		t.Errorf("Ожидался %s, получен %s", expected, result)
	}
}
