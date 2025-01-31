package main

import (
	"fmt"
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

	fmt.Println("=== Testing Visitor Pattern ===")
	services := []pattern.Service{
		&pattern.AuthService{},
		&pattern.OrderService{},
		&pattern.PaymentService{},
	}
	// Посетители
	loggingVisitor := &pattern.LoggingVisitor{}
	metricsVisitor := &pattern.MetricsVisitor{}

	// Обход всех сервисов и применение посетителей
	fmt.Println("=== Логирование запросов ===")
	for _, service := range services {
		service.Accept(loggingVisitor)
	}

	fmt.Println("\n=== Сбор метрик ===")
	for _, service := range services {
		service.Accept(metricsVisitor)
	}

}
