package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

# Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

# Дополнительное

# Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
type flags struct {
	k int
	n bool
	r bool
	u bool
}

func main() {
	f := parsLines()

	lines, err := readText("text.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sortLines, err := mySort(lines, f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, val := range sortLines {
		fmt.Println(val)
	}
}

func parsLines() flags {
	var f flags
	flag.IntVar(&f.k, "k", 1, "column number for sorting")
	flag.BoolVar(&f.n, "n", false, "sort numerically")
	flag.BoolVar(&f.r, "r", false, "sort in reverse order")
	flag.BoolVar(&f.u, "u", false, "output unique lines only")
	flag.Parse()
	f.k-- //индекс кологки с 0
	return f
}

func readText(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("файл не найдет")
	}
	defer file.Close()

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

func mySort(lines []string, f flags) ([]string, error) {
	if f.u {
		lines = repeatedLinesU(lines)
	}

	sort.Slice(lines, func(i, j int) bool {
		a := lines[i]
		b := lines[j]

		cols1, cols2 := strings.Fields(a), strings.Fields(b)

		if f.k > 0 {
			if len(cols1) >= f.k && len(cols2) >= f.k {
				a = cols1[f.k-1]
				b = cols2[f.k-1]
			} else {
				return i < j
			}
		}
		if f.n {
			numberSortN(lines)
		}
		return a < b
	})
	if f.r {
		reverseR(lines)
	}
	return lines, nil
}

// -r — сортировать в обратном порядке
func reverseR(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}

// -u — не выводить повторяющиеся строки
func repeatedLinesU(lines []string) []string {
	repetition := make(map[string]struct{}) //создаем мапу для отбора повторений
	var res []string                        //результирующий срез

	for _, val := range lines {
		if _, ok := repetition[val]; !ok { // проверяем, существует ли ключ в мапе
			repetition[val] = struct{}{}
			res = append(res, val)
		}
	}
	return res
}

// -n — сортировать по числовому значению
func numberSortN(lines []string) []string {
	// Сортируем строки с использованием числового значения
	sort.Slice(lines, func(i, j int) bool {
		num1, err1 := strconv.Atoi(lines[i]) //Преобразуем строки в числа для дальнейшего сравнения
		num2, err2 := strconv.Atoi(lines[j])

		if err1 == nil && err2 == nil {
			return num1 < num2
		}
		return lines[i] < lines[j]
	})
	return lines
}
