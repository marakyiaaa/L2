package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/*
=== Утилита wget ===
Реализовать утилиту wget с возможностью скачивать сайты целиком
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
// wget --mirror -p --convert-links -P ./<папка> адрес_сайта

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите URL: ")
	inputURL, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Ошибка чтения URL:", err)
	}
	inputURL = strings.TrimSpace(inputURL)
	if inputURL == "" {
		log.Println("URL не может быть пустым")
	}

	fmt.Print("Введите директорию для сохранения: ")
	outputDir, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Ошибка чтения директории:", err)
	}
	outputDir = strings.TrimSpace(outputDir)
	if outputDir == "" {
		log.Println("Директория не может быть пустой")
	}

	fmt.Print("Введите глубину скачивания: ")
	var depth int
	_, err = fmt.Scanln(&depth)
	if err != nil {
		log.Println("Ошибка чтения глубины:", err)
	}

	wget, err := NewWget()
	if err != nil {
		log.Println("Ошибка при создании Wget:", err)
	}

	wget.URL, err = url.Parse(inputURL)
	if err != nil {
		log.Println("Ошибка при разборе URL:", err)

	}

	wget.OutputDir = outputDir
	wget.Depth = depth

	err = wget.DownloadPage(wget.URL.String(), 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при скачивании страницы:", err)
	}

	for _, e := range wget.Errors {
		fmt.Fprintln(os.Stderr, e)
	}
}

type Wget struct {
	URL       *url.URL        //Хранение базового адреса сайта
	visit     map[string]bool //Для отслеживания посещенных URL
	Depth     int             //Для ограничения глубины рекурсивного скачивания. Это полезно, чтобы не загружать бесконечное количество страниц
	OutputDir string          //Директория для сохранения
	Errors    []error
}

func NewWget() (*Wget, error) {
	return &Wget{
		visit: make(map[string]bool),
		Depth: -1,
		URL:   &url.URL{},
	}, nil
}

func (w *Wget) DownloadPage(pageURL string, depth int) error {
	if w == nil {
		return fmt.Errorf("Wget не инициализирован")
	}

	if w.Depth != -1 && depth > w.Depth {
		return nil
	}

	if w.visit[pageURL] {
		return nil
	}
	w.visit[pageURL] = true

	resp, err := http.Get(pageURL)
	if err != nil {
		w.Errors = append(w.Errors, fmt.Errorf("ошибка при загрузке страницы %s: %w", pageURL, err))
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.Errors = append(w.Errors, fmt.Errorf("статус ответа для %s: %s", pageURL, resp.Status))
		return nil
	}

	// Сохранение страницы
	path, err := w.createPath(pageURL)
	if err != nil {
		return err
	}

	err = w.createDir(filepath.Dir(path))
	if err != nil {
		return err
	}

	var bodyCopy bytes.Buffer
	tee := io.TeeReader(resp.Body, &bodyCopy)

	err = w.saveFile(path, tee)
	if err != nil {
		return err
	}

	links := w.saveLink(&bodyCopy)
	w.processLinks(links, depth+1)

	return nil
}

// Создание пути для сохранения файла
func (w *Wget) createPath(pageURL string) (string, error) {
	parsURL, err := url.Parse(pageURL)
	if err != nil {
		w.Errors = append(w.Errors, fmt.Errorf("ошибка при разборе URL %s: %w", pageURL, err))
		return "", err
	}

	path := filepath.Join(w.OutputDir, parsURL.Path)
	if strings.HasSuffix(parsURL.Path, "/") {
		path = filepath.Join(path, "index.html")
	} else if filepath.Ext(path) == "" {
		path += ".html"
	}
	return path, nil
}

// Создание директории
func (w *Wget) createDir(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		w.Errors = append(w.Errors, fmt.Errorf("ошибка при создании директории для %s: %w", path, err))
		return err
	}
	return nil
}

// Сохранение файла
func (w *Wget) saveFile(path string, content io.Reader) error {
	file, err := os.Create(path)
	if err != nil {
		w.Errors = append(w.Errors, fmt.Errorf("ошибка при создании файла %s: %w", path, err))
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		w.Errors = append(w.Errors, fmt.Errorf("ошибка при сохранении файла %s: %w", path, err))
		return err
	}
	return nil
}

// Сохранение ссылок
func (w *Wget) saveLink(reader io.Reader) []string {
	var links []string
	z := html.NewTokenizer(reader)

	for {
		tokenT := z.Next()
		if tokenT == html.ErrorToken {
			break
		}

		token := z.Token()
		if tokenT == html.StartTagToken {
			switch token.Data {
			case "a", "link":
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			case "img", "script":
				for _, attr := range token.Attr {
					if attr.Key == "src" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
	return links
}

// Обработка ссылок
func (w *Wget) processLinks(links []string, depth int) {
	for _, link := range links {
		absoluteURL := w.resolveURL(link, w.URL.String())
		if strings.HasPrefix(absoluteURL, w.URL.String()) {
			err := w.DownloadPage(absoluteURL, depth+1)
			if err != nil {
				return
			}
		}
	}
}

func (w *Wget) resolveURL(relativeURL, baseURL string) string {
	if baseURL == "" {
		return relativeURL
	}

	u, err := url.Parse(relativeURL)
	if err != nil {
		return relativeURL
	}

	if u.IsAbs() {
		return u.String()
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return relativeURL
	}

	return base.ResolveReference(u).String()
}
