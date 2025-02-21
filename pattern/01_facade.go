package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern

Фасад — структурный паттерн проектирования уровня объекта,
который предоставляет упрощённый интерфейс к сложной подсистеме классов, библиотеке или фреймворку.

Плюсы:
 1. Снижает зависимость кода клиента от подсистемы.
 2. Изолирует клиента от компонентов подсистемы.
 3. Упрощает взаимодействие с подсистемой.

Минусы:
 1. Фасад может стать объектом, который делает слишком много/связан со всеми классами программы
и будет слишком связанный объект (рискует стать божественным объектом)
*/

// Первый объект для реализации фасада
type TV struct{}

func (tv *TV) StartTV() {
	fmt.Println("TV is started")
}

// Второй объект для реализации фасада
type Music struct {
}

func (m *Music) StartMusic() {
	fmt.Println("Music is started")
}

// Третий объект для реализации фасада
type Glow struct {
}

func (g *Glow) StartGlow() {
	fmt.Println("Glow is started")
}

// Общая структура фасада - объединения нужного функционала
type Cinema struct {
	tv    *TV
	music *Music
	glow  *Glow
}

// Конструктор для Cinema
func NewCinema() *Cinema {
	return &Cinema{tv: &TV{}, music: &Music{}, glow: &Glow{}}
}

func (c *Cinema) StartCinema() {
	c.tv.StartTV()
	c.music.StartMusic()
	c.glow.StartGlow()
}
