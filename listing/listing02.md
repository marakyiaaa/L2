Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {//2
	defer func() {
		x++
		fmt.Println(x)
	}()
	x = 1
	return
}


func anotherTest() int {//1
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest()) 
}
```

Ответ:
```
Output:
2
1

Output:
2
1

Объяснение:
- В функции test() переменная x является возвращаемым значением. Отложенная функция увеличивает x на 1 после выполнения `return`, поэтому функция возвращает 2.
- В функции anotherTest() переменная x объявлена локально. Отложенная функция увеличивает x на 1, но это изменение не влияет на возвращаемое значение, так как `return` уже выполнен. Поэтому функция возвращает 1.

```
