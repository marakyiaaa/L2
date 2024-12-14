package main

import "testing"

import (
	"reflect"
)

// Тест для функции mySort
func TestMySort(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		flags    flags
		expected []string
	}{
		{
			name:     "Сортировка без флагов",
			input:    []string{"b", "a", "c"},
			flags:    flags{},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "-r (обратный порядок)",
			input:    []string{"b", "a", "c"},
			flags:    flags{r: true},
			expected: []string{"c", "b", "a"},
		},
		{
			name:     "-u (уникальные строки)",
			input:    []string{"b", "a", "b", "c"},
			flags:    flags{u: true},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "-k (колонка для сортировки)",
			input:    []string{"2 b", "1 a", "3 c"},
			flags:    flags{k: 1},
			expected: []string{"1 a", "2 b", "3 c"},
		},
		{
			name:     "-n (числовая сортировка)",
			input:    []string{"2", "1", "3"},
			flags:    flags{n: true},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "-k и -n",
			input:    []string{"2 b", "1 a", "3 c"},
			flags:    flags{k: 1, n: true},
			expected: []string{"1 a", "2 b", "3 c"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := mySort(test.input, test.flags)
			if err != nil {
				t.Errorf("ошибка: %v", err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("ожидалось %v, получено %v", test.expected, result)
			}
		})
	}
}

func TestRepeatedLinesU(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Уникальные строки",
			input:    []string{"a", "b", "a", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Все строки уникальны",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := repeatedLinesU(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("ожидалось %v, получено %v", test.expected, result)
			}
		})
	}
}

func TestNumberSortN(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Сортировка чисел",
			input:    []string{"3", "1", "2"},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "Сортировка смешанных данных",
			input:    []string{"3", "1", "2", "a", "b"},
			expected: []string{"1", "2", "3", "a", "b"},
		},
		{
			name:     "Пустой ввод",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := numberSortN(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("ожидалось %v, получено %v", test.expected, result)
			}
		})
	}
}
