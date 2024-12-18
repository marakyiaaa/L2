package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

var testData = "apple\tbanana\tcherry\ndog\telephant\tfox\ngrape\thoney\tjuice\n"

func TestCutWithOriginal(t *testing.T) {
	tests := []struct { //структура для кейсов
		name       string
		flags      flags
		wantOutput string
	}{
		{
			name: "первый стобец",
			flags: flags{
				fields:    1,
				delimiter: "\t",
				separated: false,
			},
			wantOutput: runOriginalCut("-f1", "-d", "\t"), //используем оригинальный cut
		},
		{
			name: "разделитель - пробел и вторая колонка",
			flags: flags{
				fields:    1,
				delimiter: " ",
				separated: false,
			},
			wantOutput: runOriginalCut("-f2", "-d", " "),
		},
		{
			name: "Только разделенные линии, 3 колонка",
			flags: flags{
				fields:    3,
				delimiter: "\t",
				separated: true,
			},
			wantOutput: runOriginalCut("-f3", "-d", "\t", "-s"),
		},
		{
			name: "номер колонки вне допустимого диапазона",
			flags: flags{
				fields:    4,
				delimiter: "\t",
				separated: false,
			},
			wantOutput: "номер колонки вне допустимого диапазона",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outputBuffer bytes.Buffer
			scanner := bufio.NewScanner(strings.NewReader(testData))
			for scanner.Scan() {
				line := scanner.Text()
				output, err := myCut(line, tt.flags)
				if err != nil {
					outputBuffer.WriteString(err.Error())
					break
				} else {
					outputBuffer.WriteString(output + "\n")
				}
			}

			if outputBuffer.String() != tt.wantOutput {
				t.Errorf("myCut() output = %q, want %q", outputBuffer.String(), tt.wantOutput)
			}
		})
	}
}

// для cut
func runOriginalCut(args ...string) string {
	cmd := exec.Command("cut", args...)
	var out bytes.Buffer
	cmd.Stdin = strings.NewReader(testData)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "ошибочное в оригинале"
	}
	return out.String()
}
