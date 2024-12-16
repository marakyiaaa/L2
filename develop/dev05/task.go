package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

/*
=== Утилита grep ===
Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignoreCase-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "lineNum num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	f := parsFlags()

	//проверка на написан ли паттрен
	if flag.NArg() < 1 {
		fmt.Println("укажи pattern")
		os.Exit(1)
	}

	pattern := flag.Arg(0)

	text, err := readFile("text.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	result, count := muGrep(text, pattern, f)

	if f.count {
		fmt.Printf("%d\n", count)
	} else {
		for _, line := range result {
			fmt.Println(line)
		}
	}
}

type flags struct {
	after      int  //A
	before     int  //B
	context    int  //C
	count      bool //c
	ignoreCase bool //i
	invert     bool //v
	fixed      bool //F
	lineNum    bool //n
}

func parsFlags() flags {
	var f flags
	flag.IntVar(&f.after, "A", 0, "Print num lines of trailing context after each match")
	flag.IntVar(&f.before, "B", 0, "Print num lines of leading context before each match")
	flag.IntVar(&f.context, "C", 0, "  Print num lines of leading and trailing context surrounding each match")
	flag.BoolVar(&f.count, "c", false, "Only a count of selected lines is written to standard output")
	flag.BoolVar(&f.ignoreCase, "i", false, "Perform case insensitive matching")
	flag.BoolVar(&f.invert, "v", false, "Selected lines are those not matching any of the specified patterns")
	flag.BoolVar(&f.fixed, "F", false, "Interpret pattern as a set of fixed strings")
	flag.BoolVar(&f.lineNum, "n", false, "Print lineNum numbers")
	flag.Parse()

	if f.context > 0 {
		f.before = f.context
		f.after = f.context
	}

	return f
}

func readFile(filename string) ([]string, error) {
	var lines []string

	//открывем файл
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("файл не найден")
	}
	defer file.Close()

	//читаем файл
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.New("ошибка чтения файла")
	}
	return lines, nil
}

func muGrep(text []string, pattern string, f flags) ([]string, int) {
	var result []string
	var count int

	if f.ignoreCase {
		pattern = "(?i)" + pattern
	}

	if f.fixed {
		pattern = regexp.QuoteMeta(pattern)
	}

	re := regexp.MustCompile(pattern)

	addLines := make(map[int]bool) //для проверки добавленных строк, чтобы не было дубликаци

	// Проходка по строкам
	for i, line := range text {
		match := re.MatchString(line)

		if f.invert {
			match = !match
		}

		if match {
			count++

			if f.before > 0 || f.context > 0 {
				for j := 1; j <= f.before && i-j >= 0; j++ {
					if !addLines[i-j] { //смотрим есть ли строка
						addLineOrNum(i-j, text[i-j], &result, f)
						addLines[i-j] = true
					}
				}
			}

			// Добавляем текущую строку
			if !addLines[i] {
				addLineOrNum(i, line, &result, f)
				addLines[i] = true
			}

			if f.after > 0 || f.context > 0 {
				for j := 1; j <= f.after && i+j < len(text); j++ {
					if !addLines[i+j] {
						addLineOrNum(i+j, text[i+j], &result, f)
						addLines[i+j] = true
					}
				}
			}
		}
	}

	if f.count {
		return nil, count
	}
	return result, count
}

func addLineOrNum(i int, line string, result *[]string, f flags) {
	if f.lineNum {
		*result = append(*result, fmt.Sprintf("%d:%s", i+1, line))
	} else {
		*result = append(*result, line)
	}
}
