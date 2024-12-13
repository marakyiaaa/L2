package main

import (
	"testing"
)

func TestUnpackingString(t *testing.T) {
	//создаем структуру с поступающей строкой,и итоговой строкой,котрую мы ожидаем в ответе и bool для ошибки
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"a1b2c3", "abbccc", false},
	}

	for _, test := range tests {
		res, err := unpackingString(test.input)
		if err != nil && !test.err {
			t.Errorf("unpackingString(%q) вернула ошибку %v, ожидалось %q", test.input, err, test.expected)
		} else if res != test.expected {
			t.Errorf("unpackingString(%q) = %q, ожидалось %q", test.input, res, test.expected)
		}
	}
}
