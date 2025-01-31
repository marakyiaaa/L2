package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern

В этом паттерне мы оперируем следующими понятиями:
Command - запрос в виде объекта на выполнение;
Receiver - объект-получатель запроса, который будет обрабатывать нашу команду;
Invoker - объект-инициатор запроса.


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

//Интерфейс получателя (Receiver)
//type Command interface {
//	Execute()
//	Undo()
//}
//
//type SmartLight struct {
//	isOn bool
//}
//
//func (s *SmartLight) TurnOn() {
//	s.isOn = true
//	fmt.Println("Свет вкл")
//}
//
//func (s *SmartLight) TurnOff() {
//	s.isOn = false
//	fmt.Println("Свет выкл")
//}
//
//type LightOnCommand struct {
//	light *SmartLight
//}
//
//func (l *LightOnCommand) Execute() {
//	l.light.TurnOn()
//}
//
//func (l *LightOnCommand) Undo() {
//	l.light.TurnOff()
//}
//
//type LightOffCommand struct {
//	light *SmartLight
//}
//
//func (l *LightOffCommand) Execute() {
//	l.light.TurnOff()
//}
//
//func (l *LightOffCommand) Undo() {
//	l.light.TurnOn()
//}
//
//type ReControl struct {
//	lastCommand Command
//}
//
//func (r *ReControl) PressButton(command Command) {
//	command.Execute()
//	r.lastCommand = command
//}
//
//func (r *ReControl) PressUndo() {
//	if r.lastCommand != nil {
//		r.lastCommand.Undo()
//	}
//}

// Интерфейс получателя (Receiver)
type Device interface {
	on()
	off()
}
