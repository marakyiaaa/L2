package main

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestCompareSorts(t *testing.T) {
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
			// Тест для mySort
			mySortResult, err := mySort(test.input, test.flags)
			if err != nil {
				t.Errorf("mySort ошибка: %v", err)
			}
			if !reflect.DeepEqual(mySortResult, test.expected) {
				t.Errorf("mySort: ожидалось %v, получено %v", test.expected, mySortResult)
			}

			// Тест для sort.Slice
			sortSliceResult := test.input
			sortSlice(sortSliceResult, test.flags)
			if !reflect.DeepEqual(sortSliceResult, test.expected) {
				t.Errorf("sort.Slice: ожидалось %v, получено %v", test.expected, sortSliceResult)
			}
		})
	}
}

func sortSlice(lines []string, f flags) {
	if f.u {
		lines = repeatedLinesU(lines)
	}

	sort.Slice(lines, func(i, j int) bool {
		a := lines[i]
		b := lines[j]

		cols1 := strings.Fields(a)
		cols2 := strings.Fields(b)

		if f.k > 0 {
			if len(cols1) >= f.k && len(cols2) >= f.k {
				a = cols1[f.k-1]
				b = cols2[f.k-1]
			} else {
				return i < j
			}
		}

		if f.n {
			numA, errA := strconv.Atoi(a)
			numB, errB := strconv.Atoi(b)
			if errA == nil && errB == nil {
				return numA < numB
			}
		}

		return a < b
	})

	if f.r {
		reverseR(lines)
	}
}
