package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
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
	}, nil
}

func (w *Wget) downaldPage(pageURL string, depth int) error {
	if w.Depth != -1 && depth > w.Depth {
		return nil
	}

	if w.visit[pageURL] {
		return nil
	}
	w.visit[pageURL] = true

	// C помощью get запроса получаем содержимое страницы.
	resp, err := http.Get(pageURL)
	if err != nil {
		w.Errors = append(w.Errors, fmt.Errorf("ошибка при загрузке страницы %s: %w", pageURL, err))
		return err
	}
	defer resp.Body.Close() /////

	//// Проверяем статус ответа
	//if resp.StatusCode != http.StatusOK {
	//	w.Errors = append(w.Errors, fmt.Errorf("статус ответа для %s: %s", pageURL, resp.Status))
	//	return nil
	//}

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
	if err := os.Mkdir(filepath.Dir(path), os.ModePerm); err != nil {
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
func (w *Wget) handlingLink(links []string, URL string, dehth int) {
	for _, link := range links {
		absURL := w.resolveURL
	}
}
