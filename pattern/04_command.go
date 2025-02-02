package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern

Команда — это поведенческий паттерн проектирования, который превращает запросы в объекты,
позволяя передавать их как аргументы при вызове методов,
ставить запросы в очередь, логировать их, а также поддерживать отмену операций.

Command - запрос в виде объекта на выполнение;
Receiver - объект-получатель запроса, который будет обрабатывать нашу команду;
Invoker - объект-инициатор запроса.


Плюсы:
Убирает прямую зависимость между объектами, вызывающими операции, и объектами, которые их непосредственно выполняют.
Позволяет реализовать простую отмену и повтор операций.
Позволяет реализовать отложенный запуск операций.
Позволяет собирать сложные команды из простых.
Реализует принцип открытости/закрытости.

Минусы:
Усложняет код программы из-за введения множества дополнительных классов.
*/

// Интерфейс Command (Команда)
type Command interface {
	Execute()
	Undo()
}

// Конкретный получатель (Receiver), которое реализует интерфейс Command
type SmartLight struct {
	isOn bool
}

func (s *SmartLight) TurnOn() {
	s.isOn = true
	fmt.Println("Свет вкл")
}

func (s *SmartLight) TurnOff() {
	s.isOn = false
	fmt.Println("Свет выкл")
}

// Конкретная команда (ConcreteCommand) — Включение света
type LightOnCommand struct {
	Light *SmartLight
}

func (l *LightOnCommand) Execute() {
	l.Light.TurnOn()
}

func (l *LightOnCommand) Undo() {
	l.Light.TurnOff()
}

// Конкретная команда (ConcreteCommand) — Выключение света
type LightOffCommand struct {
	Light *SmartLight
}

func (l *LightOffCommand) Execute() {
	l.Light.TurnOff()
}

func (l *LightOffCommand) Undo() {
	l.Light.TurnOn()
}

// Invoker (Отправитель команд)
type ReControl struct {
	lastCommand Command
}

func (r *ReControl) PressButton(command Command) {
	command.Execute()
	r.lastCommand = command
}

func (r *ReControl) PressUndo() {
	if r.lastCommand != nil {
		r.lastCommand.Undo()
	}
}
