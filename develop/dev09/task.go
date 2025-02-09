package main

import (
	"bytes"
	"flag"
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
	mirror := flag.Bool("mirror", false, "Рекурсивное скачивание")
	preserve := flag.Bool("p", false, "Скачивание необходимых файлов (CSS, JS, изображения)")
	convertLinks := flag.Bool("convert-links", false, "Изменять ссылки на локальные")
	outputDir := flag.String("P", ".", "Каталог для сохранения")

	flag.Parse()

	// Проверяем, передан ли URL
	if flag.NArg() < 1 {
		fmt.Println("Использование: wget --mirror -p --convert-links -P <папка> <URL>")
		os.Exit(1)
	}
	inputURL := flag.Arg(0)

	// Парсим URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		log.Fatalf("Ошибка разбора URL: %v", err)
	}

	// Получаем абсолютный путь директории
	absOutputDir, err := filepath.Abs(*outputDir)
	if err != nil {
		log.Fatalf("Ошибка обработки пути: %v", err)
	}

	// Создаем объект Wget
	wget, err := NewWget()
	if err != nil {
		log.Fatalf("Ошибка при создании Wget: %v", err)
	}

	wget.URL = parsedURL
	wget.OutputDir = absOutputDir
	wget.Depth = -1 // поддержка глубины рекурсии

	fmt.Println("Настройки:")
	fmt.Printf("Рекурсивное скачивание: %v\n", *mirror)
	fmt.Printf("Сохранение зависимостей (CSS, JS): %v\n", *preserve)
	fmt.Printf("Конвертация ссылок: %v\n", *convertLinks)
	fmt.Printf("Каталог сохранения: %s\n", absOutputDir)

	// Запускаем скачивание
	err = wget.DownloadPage(wget.URL.String(), 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при скачивании страницы:", err)
	}

	// Выводим ошибки (если есть)
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
		return fmt.Errorf("wget не инициализирован")
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

	err = w.createDir(path)
	if err != nil {
		return err
	}

	var bodyCopy bytes.Buffer
	tee := io.TeeReader(resp.Body, &bodyCopy)

	err = w.saveFile(path, tee)
	if err != nil {
		return err
	}

	links := w.saveLink(&bodyCopy, pageURL)
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

	path := filepath.Join(w.OutputDir, parsURL.Host, parsURL.Path)
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
	log.Printf("Директория создана: %s", filepath.Dir(path))
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
func (w *Wget) saveLink(reader io.Reader, baseURL string) []string {
	var links []string
	z := html.NewTokenizer(reader)

	for {
		tokenT := z.Next()
		if tokenT == html.ErrorToken { //достигнут конец документа - цикл завершается.
			break
		}

		token := z.Token()
		if tokenT == html.StartTagToken { //Если текущий токен является открывающим тегом (<a>, <img>, <script> и т. д.), он проверяет его атрибуты.
			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					absURL := w.resolveURL(attr.Val, baseURL)
					links = append(links, absURL)
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
		if strings.HasPrefix(absoluteURL, w.URL.String()) { //принадлежит ли ссылка тому же домену
			err := w.DownloadPage(absoluteURL, depth+1) // может зациклиться, если на странице есть ссылки на уже загруженные страницы
			if err != nil {
				return
			}
		}
	}
}

// приводит link (которая может быть относительной) к абсолютному UR
func (w *Wget) resolveURL(relativeURL, baseURL string) string {
	if baseURL == "" {
		return relativeURL
	}

	u, err := url.Parse(relativeURL)
	if err != nil {
		return relativeURL
	}

	// является ли абсолютным
	if u.IsAbs() {
		return u.String()
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return relativeURL
	}

	return base.ResolveReference(u).String()
}
