package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	f := parsFlags()
	readStdin(f)

}

type flags struct {
	fields    int
	delimiter string
	separated bool
}

func readStdin(f flags) {
	if flag.NArg() == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			output, err := myCut(line, f)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			fmt.Println(output)
		}
	}
}

func myCut(str string, f flags) (string, error) {
	columns := strings.Split(str, f.delimiter) //разбиваю str по разделителю \t

	//f и номер колонки в допустимом диапазоне
	if f.fields > 0 {
		if f.fields > len(columns) {
			return "", errors.New("номер колонки вне допустимого диапазона")
		}
		return columns[f.fields-1], nil
	}

	//s и строка не содержит разделителя
	if f.separated && !strings.Contains(str, f.delimiter) {
		return "", nil
	}

	//не -f - возвращаем всю строку
	return strings.Join(columns, f.delimiter), nil

}

func parsFlags() flags {
	var f flags
	flag.IntVar(&f.fields, "f", 0, "выбрать поля (колонки)")
	flag.StringVar(&f.delimiter, "d", "\t", "использовать другой разделитель")
	flag.BoolVar(&f.separated, "s", false, "только строки с разделителем")
	flag.Parse()
	return f
}
