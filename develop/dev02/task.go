package main

import (
	"fmt"
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
	fmt.Println(unpackingString("222"))
}

func unpackingString(str string) (string, error) {
	runes := []rune(str)
	var res []rune
	var num int32 = 0
	var count int32 = 1
	var j int32

	for i := len(runes) - 1; i >= 0; i-- {
		if !unicode.IsNumber(runes[i]) { //если значене не цифра
			for j = 0; j < num-1; j++ {
				res = append(res, runes[i])
			}
			res = append(res, runes[i])
			num = 0
			count = 1
		} else if unicode.IsNumber(runes[i]) {
			num += int32(runes[i]) * count
			count *= 10
			fmt.Println(num, count)
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
