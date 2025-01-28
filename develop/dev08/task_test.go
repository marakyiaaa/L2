package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestCD(t *testing.T) {
	tests := []struct {
		args    []string
		wantErr bool
	}{
		{[]string{"cd", ".."}, false},
		{[]string{"cd", "/nonexistent"}, true},
	}

	for _, tt := range tests {
		err := cd(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("cd(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
		}
	}
}

func TestPWD(t *testing.T) {
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Ошибка получения текущей директории: %v", err)
	}
	defer os.Chdir(oldDir)

	err = pwd()
	if err != nil {
		t.Errorf("pwd() error = %v, wantErr false", err)
	}
}

func TestEcho(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"echo", "hello", "world"}, "hello world\n"},
	}

	for _, tt := range tests {
		output := captureOutput(func() {
			echo(tt.args)
		})
		if output != tt.want {
			t.Errorf("echo(%v) = %q, want %q", tt.args, output, tt.want)
		}
	}
}

func TestKill(t *testing.T) {
	cmd := exec.Command("sleep", "100")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Ошибка запуска процесса: %v", err)
	}
	defer cmd.Process.Kill()

	err := kill([]string{"kill", strconv.Itoa(cmd.Process.Pid)})
	if err != nil {
		t.Errorf("kill() error = %v, wantErr false", err)
	}
}

func TestPS(t *testing.T) {
	err := ps()
	if err != nil {
		t.Errorf("ps() error = %v, wantErr false", err)
	}
}

func TestPipe(t *testing.T) {
	tests := []struct {
		commands []string
		want     string
	}{
		{[]string{"echo hello", "wc -c"}, "       6\n"},
		{[]string{"echo test", "wc -c"}, "       5\n"},
	}

	for _, tt := range tests {
		output := captureOutput(func() {
			err := pipe(tt.commands)
			if err != nil {
				t.Errorf("pipe(%v) error = %v", tt.commands, err)
			}
		})
		if output != tt.want {
			t.Errorf("pipe(%v) = %q, want %q", tt.commands, output, tt.want)
		}
	}
}
func TestNC(t *testing.T) {
	go func() {
		ln, _ := net.Listen("tcp", "localhost:30308")
		defer ln.Close()

		conn, _ := ln.Accept()
		defer conn.Close()

		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		if string(buf[:n]) != "test\n" {
			t.Errorf("Ожидаемые данные: %q, получено: %q", "test\n", string(buf[:n]))
		}
	}()

	err := nc([]string{"nc", "localhost", "30308"})
	if err != nil {
		t.Errorf("nc() error = %v, wantErr false", err)
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
