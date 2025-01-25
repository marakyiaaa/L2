package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">>> ")
		input, err := readInput(reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if input == "" {
			continue
		}

		commands := strings.Split(input, "|")
		if len(commands) > 1 {
			if err := pipe(commands); err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка:", err)
			}
			continue
		}

		args := strings.Fields(input)
		if err := handlesCommand(args); err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка:", err)
		}
	}
}

func readInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ввода: %v", err)
	}
	return strings.TrimSpace(input), nil
}

func handlesCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("не указана команда")
	}

	// в фоновом режиме
	background := false
	if args[len(args)-1] == "&" {
		background = true
		args = args[:len(args)-1]
	}

	switch args[0] {
	case "cd":
		return cd(args)
	case "pwd":
		return pwd()
	case "echo":
		return echo(args)
	case "kill":
		return kill(args)
	case "ps":
		return ps()
	default:
		return execute(args, background)
	}
}

func execute(args []string, background bool) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if background {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("ошибка запуска команды: %v", err)
		}
		fmt.Printf("[%d] Запущено в фоновом режиме\n", cmd.Process.Pid)
		return nil
	}

	// Запуск и ожидание завершения
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка выполнения команды: %v", err)
	}
	return nil
}

func cd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("cd: недостаточно аргументов")
	}

	dir := args[1]
	if dir == ".." {
		if err := os.Chdir(".."); err != nil {
			return fmt.Errorf("cd ..: %v", err)
		}
		return nil
	}

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("cd: %v", err)
	}
	return nil
}

func pwd() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("pwd: %v", err)
	}
	fmt.Println(dir)
	return nil
}

func echo(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("echo: недостаточно аргументов")
	}
	fmt.Println(strings.Join(args[1:], " "))
	return nil
}

func kill(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("kill: недостаточно аргументов")
	}

	pid, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("kill: неверный PID: %v", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("kill: %v", err)
	}

	if err := process.Signal(syscall.SIGKILL); err != nil {
		return fmt.Errorf("kill: %v", err)
	}
	return nil
}

func ps() error {
	cmd := exec.Command("ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ps: %v", err)
	}
	return nil
}

func pipe(commands []string) error {
	var predcmd *exec.Cmd        // cсылка на предыдущую команду
	var prevStdout io.ReadCloser // пайп для вывода предыдущей команды

	for i, cmdstr := range commands {
		args := strings.Fields(strings.TrimSpace(cmdstr))
		if len(args) == 0 {
			return fmt.Errorf("пустая команда")
		}

		cmd := exec.Command(args[0], args[1:]...)

		// если это не первая команда, связываем её ввод с выводом предыдущей команды
		if i > 0 {
			cmd.Stdin = prevStdout
		}

		// если это не последняя команда, создаем пайп для её вывода
		if i < len(commands)-1 {
			var err error
			prevStdout, err = cmd.StdoutPipe()
			if err != nil {
				return fmt.Errorf("ошибка создания пайпа: %v", err)
			}
		} else {
			cmd.Stdout = os.Stdout
		}

		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("ошибка выполнения команды: %v", err)
		}

		predcmd = cmd
	}

	if err := predcmd.Wait(); err != nil {
		return fmt.Errorf("ошибка выполнения команды: %v", err)
	}

	return nil
}
