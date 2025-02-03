package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern

Стратегия — это поведенческий паттерн проектирования, который определяет семейство схожих алгоритмов и помещает каждый из них в собственный класс,
после чего алгоритмы можно взаимозаменять прямо во время исполнения программы.

Плюсы:
Быстрая замена алгоритмов.
Изолирует код и данные алгоритмов от остальных классов.
Уход от наследования к делегированию.
Реализует принцип открытости/закрытости.

Минусы:
Усложняет программу за счёт дополнительных классов.
Клиент должен знать, в чём состоит разница между стратегиями, чтобы выбрать подходящую.

В паттерне "Стратегия" есть:
Контекст (Context) — объект, который использует стратегию.
Стратегия (Strategy) — интерфейс, который определяет поведение.
Конкретная стратегия (ConcreteStrategy) — различные реализации этого интерфейса,
которые могут изменяться в зависимости от ситуации.
*/

// Strategy -  Интерфейс, который будет использоваться в контексте (классе Cache).
type EvictionAlgo interface {
	evict(c *Cache)
}

// Context
type Cache struct {
	storage      map[string]string
	evictionAlgo EvictionAlgo
	capacity     int
	maxCapacity  int
}

// Конкретные стратегии
type Fifo struct{}

func (f *Fifo) evict(c *Cache) {
	fmt.Println("Fifo strategy")
}

// Конкретные стратегии
type Lifo struct{}

func (l *Lifo) evict(c *Cache) {
	fmt.Println("Lifo strategy")
}

// Позволяет изменять стратегию вытеснения в любой момент времени.
// Это позволяет менять поведение кэша без изменения его кода, просто выбирая другую стратегию.
func (c *Cache) SetEvict(e EvictionAlgo) {
	c.evictionAlgo = e
}

func (c *Cache) evict() {
	if c.evictionAlgo != nil {
		c.evictionAlgo.evict(c)
		c.capacity--
	} else {
		fmt.Println("Error")
	}
}

func (c *Cache) Add(key, value string) {
	if c.capacity == c.maxCapacity {
		c.evict()
	}
	c.capacity++
	c.storage[key] = value
}

func (c *Cache) Get(key string) string {
	value, exists := c.storage[key]
	if exists {
		return value
	}
	return ""
}

func InitCache(e EvictionAlgo) *Cache {
	storage := make(map[string]string)
	return &Cache{
		storage:      storage,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}
