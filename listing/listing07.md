Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
asChan принимает значения из среза и записываетв канал, после канал закрывается
так же между отправками значений происходит рандомная задержка

merge соединяет 2 канала в 1
черех select считываются данные с a,b и записываются в с
но нет провнрки на закрыте каналов и из за этого,когда каналы закрываются
прогамма получает бесконечно 0

вывод - случайный порядо чисел из слайсов и далее бесконечное количество 0
```
