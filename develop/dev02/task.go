package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

/*
=== Задача на распаковку ===
Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)
В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.
Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	fmt.Println(unpackingString("a4bc2d5e"))
}

func unpackingString(str string) (string, error) {
	runes := []rune(str)
	var res []rune
	num := 0
	count := 1

	if len(runes) > 0 && unicode.IsDigit(runes[0]) {
		return "", errors.New("некорректная строка")
	}

	for i := len(runes) - 1; i >= 0; i-- { //идем с конца <-
		if !unicode.IsNumber(runes[i]) { //если значе6ние не число
			for j := 0; j < num-1; j++ {
				res = append(res, runes[i])
			}
			res = append(res, runes[i])
			num = 0
			count = 1
		} else if unicode.IsNumber(runes[i]) {
			convStr := string(runes[i])
			tmp, _ := strconv.Atoi(convStr)
			num += tmp * count
			count *= 10
		}
	}
	return reverseString(res), nil
}

func reverseString(runes []rune) string {
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
