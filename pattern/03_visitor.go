package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern

Посетитель — это поведенческий паттерн проектирования,
который позволяет добавлять в программу новые операции,
не изменяя классы объектов, над которыми эти операции могут выполняться.

Accept() нужен для передачи управления от объекта к посетителю.

Двойная диспетчеризация помогает вызвать правильный метод для конкретного типа.

Плюсы:
Упрощает добавление операций, работающих со сложными структурами объектов.
Объединяет родственные операции в одном классе.
Посетитель может накапливать состояние при обходе структуры элементов.

Минусы:
Паттерн не оправдан, если иерархия элементов часто меняется.
Может привести к нарушению инкапсуляции элементов.
*/

// интерфейс всех микросервисов для примера
type Service interface {
	Accept(v Visitor) // Метод для принятия посетителя
}

// интерфейс для всех посетителей
type Visitor interface {
	VisitAuthService(*AuthService)
	VisitOrderService(*OrderService)
	VisitPaymentService(*PaymentService)
}

// сервис аутентификации
type AuthService struct{}

func (a *AuthService) Accept(v Visitor) {
	v.VisitAuthService(a) // Передаёт себя посетителю
}

// сервис заказов
type OrderService struct{}

func (o *OrderService) Accept(v Visitor) {
	v.VisitOrderService(o) // Передаёт себя посетителю
}

// сервис платежей
type PaymentService struct{}

func (p *PaymentService) Accept(v Visitor) {
	v.VisitPaymentService(p) // Передаёт себя посетителю
}

// Реализация Посетителя для логирования
type LoggingVisitor struct{}

func (l *LoggingVisitor) VisitAuthService(a *AuthService) {
	fmt.Println("[LOG] Запрос в AuthService")
}

func (l *LoggingVisitor) VisitOrderService(o *OrderService) {
	fmt.Println("[LOG] Запрос в OrderService")
}

func (l *LoggingVisitor) VisitPaymentService(p *PaymentService) {
	fmt.Println("[LOG] Запрос в PaymentService")
}

// Реализация Посетителя для метрик
type MetricsVisitor struct{}

func (m *MetricsVisitor) VisitAuthService(a *AuthService) {
	fmt.Println("[METRICS] +1 запрос в AuthService")
}

func (m *MetricsVisitor) VisitOrderService(o *OrderService) {
	fmt.Println("[METRICS] +1 запрос в OrderService")
}

func (m *MetricsVisitor) VisitPaymentService(p *PaymentService) {
	fmt.Println("[METRICS] +1 запрос в PaymentService")
}
