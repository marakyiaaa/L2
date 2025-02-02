package main

import (
	"fmt"
	"math/rand"
	"pattern/pattern"
)

func main() {
	//fmt.Println("=== Testing Facade Pattern ===")
	//cinema := pattern.NewCinema()
	//cinema.StartCinema()
	//
	//fmt.Println("=== Testing Builder Pattern ===")
	//builder := pattern.NewConcreteBuilder()
	//director := pattern.NewJSONDirector(builder)
	//jsonString, _ := director.BuildSampleJSON()
	//fmt.Println(jsonString)
	//
	//fmt.Println("=== Testing Builder Pattern (PersonBuilder) ===")
	//person := pattern.NewPersonBuilder().
	//	SetName("Alice").
	//	SetAge(30).
	//	SetAddress("Moscow", "Russia")
	//jsonPerson, _ := person.GetJSON()
	//fmt.Println("JSON для человека:", jsonPerson)

	//fmt.Println("=== Testing Visitor Pattern ===")
	//services := []pattern.Service{
	//	&pattern.AuthService{},
	//	&pattern.OrderService{},
	//	&pattern.PaymentService{},
	//}
	//// Посетители
	//loggingVisitor := &pattern.LoggingVisitor{}
	//metricsVisitor := &pattern.MetricsVisitor{}
	//
	//// Обход всех сервисов и применение посетителей
	//fmt.Println("=== Логирование запросов ===")
	//for _, service := range services {
	//	service.Accept(loggingVisitor)
	//}
	//
	//fmt.Println("\n=== Сбор метрик ===")
	//for _, service := range services {
	//	service.Accept(metricsVisitor)
	//}

	//fmt.Println("=== Testing Command Pattern ===")
	//light := &pattern.SmartLight{}
	//controller := &pattern.ReControl{}
	//lightOn := &pattern.LightOnCommand{Light: light}
	//lightOff := &pattern.LightOffCommand{Light: light}
	//
	//// Нажимаем кнопку включения света
	//controller.PressButton(lightOn)
	//// Нажимаем кнопку выключения света
	//controller.PressButton(lightOff)
	//// Отмена последней команды (включаем свет обратно)
	//controller.PressUndo()

	//fmt.Println("=== Testing Chain Of Responsibility Pattern ===")
	//spamFilter := &pattern.SpamFilter{}
	//authentication := &pattern.Authentication{}
	//access := &pattern.Authorization{}
	//server := &pattern.Server{}
	//
	//// Строим цепочку: спам → аутентификация → авторизация → сервер
	//spamFilter.SetNext(authentication).SetNext(access).SetNext(server)
	//
	//fmt.Println("\n --- Администратор отправляет сообщение --- \n")
	//request1 := pattern.Request{User: "Admin", Status: "admin", Message: "Привет!", IsLoggedIn: true}
	//spamFilter.Handle(request1)
	//
	//fmt.Println(" \n --- Гость без аутентификации --- \n")
	//request2 := pattern.Request{User: "Guest", Status: "guest", Message: "Привет!", IsLoggedIn: false}
	//spamFilter.Handle(request2)

	fmt.Println("=== Testing Factory Method Pattern ===")
	providers := []string{"paypal", "visa", "master"}
	randomProvider := providers[rand.Intn(len(providers))]

	paymentGet := pattern.GetPaymant(randomProvider)

	if paymentGet != nil {
		fmt.Println("Выбран провайдер:", randomProvider)
		fmt.Println(paymentGet.ProcessPayment(10010101.50))
	} else {
		fmt.Println("Ошибка: Неизвестная платёжная система -", randomProvider)
	}
}
