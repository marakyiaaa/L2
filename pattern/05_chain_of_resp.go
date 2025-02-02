package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

Цепочка обязанностей — это поведенческий паттерн проектирования,
который позволяет передавать запросы последовательно по цепочке обработчиков.
Каждый последующий обработчик решает,
может ли он обработать запрос сам и стоит ли передавать запрос дальше по цепи.

Chain Of Responsibility позволяет избежать привязки объекта-отправителя запроса
к объекту-получателю запроса, при этом давая шанс обработать этот запрос нескольким объектам.
Получатели связываются в цепочку, и запрос передается по цепочке,
пока не будет обработан каким-то объектом.

Вместо хранения ссылок на всех кандидатов-получателей запроса,
каждый отправитель хранит единственную ссылку на начало цепочки,
а каждый получатель имеет единственную ссылку на своего преемника - последующий элемент в цепочке.

Плюсы:
Уменьшает зависимость между клиентом и обработчиками.
Реализует принцип единственной обязанности.
Реализует принцип открытости/закрытости.

Минусы:
Запрос может остаться никем не обработанным.
*/

/*
Цепочка обработчиков, которая проверяет:

Фильтр спама → блокирует сообщения с запрещёнными словами.
Аутентификацию → проверяет, вошёл ли пользователь.
Авторизацию → проверяет, есть ли у пользователя права доступа.

Если сообщение проходит все проверки, оно отправляется на сервер.
Если не проходит хоть одну → цепочка прерывается.
*/

// Интерфейс обработчика (Handler)
type Handler interface {
	SetNext(handler Handler) Handler // Устанавливает следующий обработчик
	Handle(request Request) bool     // Обрабатывает запрос
}

//Каждый обработчик должен уметь:
//Принимать запрос (Handle)
//Передавать его дальше (SetNext)

type Request struct {
	User       string
	Status     string //(admin, user, guest)
	Message    string
	IsLoggedIn bool // Вошёл ли пользователь в систему
}

//Базовый обработчик (чтобы не дублировать код)
//BaseHandler передаёт запрос дальше по цепочке.

type BaseHandler struct {
	next Handler
}

// Установка следующего обработчика
func (bh *BaseHandler) SetNext(handler Handler) Handler {
	bh.next = handler
	return handler // Позволяет строить цепочку
}

// Передаёт запрос дальше, если есть следующий обработчик
func (bh *BaseHandler) Handle(request Request) bool {
	if bh.next != nil {
		return bh.next.Handle(request)
	}
	return true // Если обработчиков больше нет, пропускаем запрос
}

// Обработичк - фильтр спама
type SpamFilter struct {
	BaseHandler
}

func (s *SpamFilter) Handle(request Request) bool {
	if request.Message == "спам" {
		fmt.Println("Сообщение - спам")
		return false // Останавливаем цепочку
	}
	fmt.Println("Сообщение проверено на спам")
	return s.BaseHandler.Handle(request) //Передаем по цепочке далее
}

// Обработичк - аутентификации
type Authentication struct {
	BaseHandler
}

func (a *Authentication) Handle(request Request) bool {
	if !request.IsLoggedIn {
		fmt.Println("Пользователь не аутентифицирован")
		return false
	}
	fmt.Println("Пользователь аутентифицирован")
	return a.BaseHandler.Handle(request) // Передаём по цепочке
}

// Обработичк -  авторизации
type Authorization struct {
	BaseHandler
}

func (a *Authorization) Handle(request Request) bool {
	if request.Status != "admin" {
		fmt.Println("Нет прав доступа")
		return false // Останавливаем цепочку
	}
	fmt.Println("Доступ разрешён")
	return a.BaseHandler.Handle(request) //Передаем по цепочке далее
}

// Обработчик - сервер
type Server struct {
	BaseHandler
}

func (s *Server) Handle(request Request) bool {
	fmt.Println("Сообщенеи отправлено на сервер")
	return true
}
