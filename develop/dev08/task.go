package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
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
	case "nc":
		return nc(args)
	default:
		return execute(args)
	}
}

func execute(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

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

func nc(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("недостаточно аргументов: nc <host> <port> [udp]")
	}

	host := args[1]
	port := args[2]
	protocol := "tcp"
	if len(args) > 3 && strings.ToLower(args[3]) == "udp" {
		protocol = "udp"
	}

	conn, err := net.Dial(protocol, fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return fmt.Errorf("ошибка подключения к %s:%s: %v", host, port, err)
	}
	defer conn.Close()

	fmt.Printf("Подключено к %s:%s (%s)\n", host, port, protocol)

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Printf("Ошибка при чтении данных: %v\n", err)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, err := conn.Write(scanner.Bytes())
		if err != nil {
			return fmt.Errorf("ошибка при отправке данных: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка при чтении stdin: %v", err)
	}
	return nil
}
