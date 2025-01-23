Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
Тип error — это интерфейс:
type error interface {
    Error() string
}

В func Foo переменная err типа *os.PathError == nil
Возварщаемое значение тип error -> создается интрфейс, где
тип  - *os.PathError
значение - nil


1) fmt.Println(err) выводится именно значение интерфейса - nil

2) fmt.Println(err == nil) здесь мы сравниваем именно интерфейс с nil,
а так как там есть поле типа данный (в данном случе *os.PathError),
то он уже не nil - выведется false.

Интерфейс считается nil только если и его ТИП nil, и ЗНАЧЕНИЕ  nil.
```
