package main

import (
	"bufio"
	"errors"
	"os"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {

}

type flags struct {
	k int
	n bool
	r bool
	u bool
}

func readText(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("файл не найдет")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.New("ошибка чтения файла")
	}
	return lines, nil
}

//func mySort(str string) (string, error) {
//
//}
