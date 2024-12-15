package main

import (
	"reflect"
	"testing"
)

func TestVocabulary(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name:     "Пустой ввод",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:     "Одно слово",
			input:    []string{"пятак"},
			expected: map[string][]string{},
		},
		{
			name:  "Дубликаты",
			input: []string{"пятак", "пятак", "пятка", "тяпка", "тяпка"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:     "Слова без анаграмм",
			input:    []string{"пятак", "кот", "собака", "дом"},
			expected: map[string][]string{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := vocabulary(&test.input)
			if !reflect.DeepEqual(*result, test.expected) {
				t.Errorf("Ожидалось %v, получено %v", test.expected, *result)
			}
		})
	}
}
