package pattern

import "encoding/json"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern

Строитель — это порождающий паттерн проектирования уровня объекта,
который позволяет создавать сложные объекты пошагово.
Строитель даёт возможность использовать один и тот же код строительства
для получения разных представлений объектов.

Плюсы:
Позволяет использовать один и тот же код для создания различных продуктов.
Изолирует сложный код сборки продукта от его основной бизнес-логики.
Позволяет собирать объекты пошагово, вызывая только те шаги, которые вам нужны.

Минусы:
Клиент будет привязан к конкретным классам строителей,
так как в интерфейсе директора может не быть метода получения результата.

После того как будет построена последняя его часть, продукт можно использовать.
*/

// JSONBuilder интерфейс строителя.
type JSONBuilder interface {
	AddString(key, value string) JSONBuilder
	AddNumber(key string, value float64) JSONBuilder
	AddBool(key string, value bool) JSONBuilder
	AddArray(key string, value []interface{}) JSONBuilder
	AddObj(key string, value map[string]interface{}) JSONBuilder
	GetJSON() (string, error)
}

//Это означает: "этот метод должен возвращать что-то, что реализует JSONBuilder".
//В ConcreteBuilder этот метод возвращает b, а b – это *ConcreteBuilder, который реализует JSONBuilder.

// конкретный строитель - ConcreteBuilder реализует интерфейс JSONBuilder.
type ConcreteBuilder struct {
	data map[string]interface{}
}

// Конструктор ConcreteBuilder.
func NewConcreteBuilder() *ConcreteBuilder {
	return &ConcreteBuilder{data: make(map[string]interface{})}
}

// Реализация функционала
func (b *ConcreteBuilder) AddString(key, value string) JSONBuilder {
	b.data[key] = value
	return b
}

func (b *ConcreteBuilder) AddNumber(key string, value float64) JSONBuilder {
	b.data[key] = value
	return b
}

func (b *ConcreteBuilder) AddBool(key string, value bool) JSONBuilder {
	b.data[key] = value
	return b
}

func (b *ConcreteBuilder) AddArray(key string, value []interface{}) JSONBuilder {
	b.data[key] = value
	return b
}

func (b *ConcreteBuilder) AddObj(key string, value map[string]interface{}) JSONBuilder {
	b.data[key] = value
	return b
}

func (b *ConcreteBuilder) GetJSON() (string, error) {
	data, err := json.MarshalIndent(b.data, "", " ") //Go-структуры в JSON.
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Директор, управляющий созданием JSON
type JSONDirector struct {
	builder JSONBuilder
}

func NewJSONDirector(builder JSONBuilder) *JSONDirector {
	return &JSONDirector{builder: builder}
}

func (d *JSONDirector) BuildSampleJSON() (string, error) {
	return d.builder.
		AddString("name", "Kate").
		AddNumber("age", 25).
		AddBool("yes", true).
		AddArray("pets", []interface{}{"cat", "dog", "fish"}).
		AddObj("address", map[string]interface{}{
			"city":    "Kazan",
			"country": "Russia",
		}).
		GetJSON()
}

// Новый строитель PersonBuilder
type PersonBuilder struct {
	ConcreteBuilder
}

// Конструктор PersonBuilder
func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{*NewConcreteBuilder()}
}

func (b *PersonBuilder) SetName(name string) *PersonBuilder {
	b.AddString("name", name)
	return b
}

func (b *PersonBuilder) SetAge(age int) *PersonBuilder {
	b.AddNumber("age", float64(age))
	return b
}

func (b *PersonBuilder) SetAddress(city, country string) *PersonBuilder {
	b.AddObj("address", map[string]interface{}{
		"city":    city,
		"country": country,
	})
	return b
}
