package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

func main() {

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		//sig(2*time.Hour),
		//sig(5*time.Minute),
		sig(1 * time.Second),
		//sig(1*time.Hour),
		//sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}

func or(channels ...<-chan any) <-chan any { // использую any - т.к. type any = interface{}
	if len(channels) == 0 {
		return nil
	}

	done := make(chan any) //результирующий ранал
	var g errgroup.Group   //errgroup для управления горутинами
	var mu sync.Once       //для синхронизации закрытия канала - гарантирует выполнение операциии ОДИН раз

	for _, ch := range channels {
		ch := ch            //для замыкания
		g.Go(func() error { //метод Go-
			for range ch {
				//читаем из канала до его закрытия
			}
			mu.Do(func() {
				close(done)
			})
			return nil
		})
	}

	go func() { //ожидание завершения горутин
		err := g.Wait()
		if err != nil {
			return
		}
		mu.Do(func() {
			close(done)
		})
	}()
	return done
}

/*
//select
func or2(channels ...<-chan any) <-chan any {
	if len(channels) == 0 {
		return nil
	}

	done := make(chan any)

	for _, val := range channels {
		go func(ch <-chan any) {
			select {
			case _, open := <-ch:
				if !open {
					close(done)
				}
			}
		}(val)
	}
	return done
}

//WaitGroup
func or3(channels ...<-chan any) <-chan any {
	if len(channels) == 0 {
		return nil
	}

	done := make(chan any)
	wg := sync.WaitGroup{}

	wg.Add(len(channels))

	for _, ch := range channels {
		go func(ch <-chan any) {
			defer wg.Done()
			for c := range ch {
				done <- c
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}*/
