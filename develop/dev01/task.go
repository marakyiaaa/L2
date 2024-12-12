package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"time"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время/точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	t, err := myTime()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Текущее время:", t)
}

func myTime() (time.Time, error) {

	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return time.Time{}, fmt.Errorf("Ошибка вывода времени")
	}
	return t, err
}

/*
Network Time Security (NTS) - используется для получения точного времени с серверов NTP
*/
