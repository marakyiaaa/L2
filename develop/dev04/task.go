package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.
Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества.
Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	str := []string{"пятак", "пЯтка", "тяпка", "листОк", "слитоК", "столик"}

	for key, words := range *vocabulary(&str) {
		fmt.Println(key, words)
	}
}

func vocabulary(words *[]string) *map[string][]string {
	anagrams := make(map[string][]string) // Мапа для хранения анаграмм
	unique := make(map[string]bool)       // Мапа для отслеживания уникальных слов

	for _, word := range *words {
		runes := []rune(strings.ToLower(word)) // Преобразуем строку в срез рун и нижний регистр
		slices.Sort(runes)                     // Сортируем

		sortedWords := string(runes) // Преобразуем обратно в строку

		if _, ok := anagrams[sortedWords]; !ok {
			//создаем новое множество анаграмм
			anagrams[sortedWords] = []string{word}
		} else {
			// Если ключ есть, добавляем слово в множество, если оно еще не было добавлено
			if !unique[word] {
				anagrams[sortedWords] = append(anagrams[sortedWords], word)
			}
		}
		unique[word] = true
	}

	// Убираем множества из одного элемента
	result := make(map[string][]string)
	for _, words := range anagrams {
		if len(words) > 1 {
			sort.Strings(words)
			result[words[0]] = words
		}
	}
	return &result
}
