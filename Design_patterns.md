Больше информации - https://refactoringguru.cn/ru/design-patterns/catalog

# Шаблоны проектирования

## Введение

Шаблоны проектирования (Design Patterns) — это проверенные временем решения типичных проблем, с которыми сталкиваются разработчики при проектировании программного обеспечения. Они помогают сделать код более читаемым, масштабируемым и поддерживаемым. Впервые эти концепции были систематизированы в книге "Design Patterns: Elements of Reusable Object-Oriented Software" (авторы: Эрих Гамма и др.), известной как "книга банды четырёх" (Gang of Four, GoF).

Шаблоны делятся на три основные категории:

- **Порождающие (Creational)** — создание объектов.
- **Структурные (Structural)** — организация классов и объектов.
- **Поведенческие (Behavioral)** — взаимодействие между объектами.

В этой лекции мы разберём ключевые идеи шаблонов и реализуем 10 самых популярных из них на языке Go, который, несмотря на свою простоту, требует особого подхода к шаблонам из-за отсутствия классического ООП (например, наследования).

## Почему шаблоны важны для Senior-разработчиков?

- **Гибкость**: Шаблоны позволяют легко адаптировать код к новым требованиям.
- **Коммуникация**: Это универсальный язык, понятный всем разработчикам.
- **Масштабируемость**: Помогают избежать "спагетти-кода" в больших проектах.
- **Поддерживаемость**: Упрощают рефакторинг и тестирование.

Однако важно помнить: шаблоны — это не "серебряная пуля". Используйте их только там, где они действительно решают проблему, чтобы не усложнять код без необходимости.

## Особенности Go и шаблонов

Go — это язык с минималистичным синтаксисом, без классов и наследования. Вместо этого он использует:

- **Интерфейсы** (implicit implementation).
- **Встраивание** (embedding).
- **Композицию** вместо наследования.

Из-за этого некоторые классические шаблоны (например, с использованием наследования) в Go реализуются иначе. Мы будем адаптировать их под идиомы Go, такие как "принимай интерфейсы, возвращай структуры".

## 10 популярных шаблонов с примерами на Go

### 1. Singleton (Одиночка)
**Назначение:** Гарантирует, что у класса есть только один экземпляр, и предоставляет глобальную точку доступа к нему.
**Пример:** Логгер, который должен быть одним для всей системы.
```go
package main

import (
	"fmt"
	"sync"
)

type Logger struct {
	logs []string
}

var instance *Logger
var once sync.Once

func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{}
	})
	return instance
}

func (l *Logger) Log(message string) {
	l.logs = append(l.logs, message)
	fmt.Println("Log:", message)
}

func main() {
	logger1 := GetLogger()
	logger2 := GetLogger()

	logger1.Log("First message")
	logger2.Log("Second message")

	fmt.Println(logger1 == logger2) // true
	fmt.Println(logger1.logs)      // [First message Second message]
}
```
**Особенности в Go:** Используем sync.Once для потокобезоп

### 2. Factory Method (Фабричный метод)
**Назначение:** Определяет интерфейс для создания объекта, но позволяет подклассам решать, какой класс инстанцировать.
**Пример:** Создание разных типов транспорта.
```go
package main

import "fmt"

type Vehicle interface {
	Drive() string
}

type Car struct{}

func (c *Car) Drive() string { return "Car driving" }

type Bike struct{}

func (b *Bike) Drive() string { return "Bike driving" }

func NewVehicle(vType string) Vehicle {
	switch vType {
	case "car":
		return &Car{}
	case "bike":
		return &Bike{}
	default:
		return nil
	}
}

func main() {
	car := NewVehicle("car")
	bike := NewVehicle("bike")

	fmt.Println(car.Drive())  // Car driving
	fmt.Println(bike.Drive()) // Bike driving
}
```
Особенности в Go: Используем функцию вместо метода, так как в Go нет наследования.

### 3. Abstract Factory (Абстрактная фабрика)
**Назначение:** Предоставляет интерфейс для создания семейств связанных объектов без указания их конкретных классов.
**Пример:** Фабрика для создания UI-компонентов (кнопки и чекбоксы).
```go
package main

import "fmt"

type Button interface {
	Click() string
}

type Checkbox interface {
	Check() string
}

type WinButton struct{}
func (w *WinButton) Click() string { return "Windows Button clicked" }

type MacButton struct{}
func (m *MacButton) Click() string { return "Mac Button clicked" }

type WinCheckbox struct{}
func (w *WinCheckbox) Check() string { return "Windows Checkbox checked" }

type MacCheckbox struct{}
func (m *MacCheckbox) Check() string { return "Mac Checkbox checked" }

type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

type WinFactory struct{}
func (w *WinFactory) CreateButton() Button     { return &WinButton{} }
func (w *WinFactory) CreateCheckbox() Checkbox { return &WinCheckbox{} }

type MacFactory struct{}
func (m *MacFactory) CreateButton() Button     { return &MacButton{} }
func (m *MacFactory) CreateCheckbox() Checkbox { return &MacCheckbox{} }

func NewGUIFactory(os string) GUIFactory {
	if os == "windows" {
		return &WinFactory{}
	}
	return &MacFactory{}
}

func main() {
	factory := NewGUIFactory("mac")
	btn := factory.CreateButton()
	chk := factory.CreateCheckbox()

	fmt.Println(btn.Click())  // Mac Button clicked
	fmt.Println(chk.Check()) // Mac Checkbox checked
}
```

### 4. Builder (Строитель)
**Назначение:** Отделяет конструирование сложного объекта от его представления.
**Пример:** Постройка дома с разными конфигурациями.
```go
package main

import "fmt"

type House struct {
	walls, doors, windows int
}

type HouseBuilder struct {
	house *House
}

func NewHouseBuilder() *HouseBuilder {
	return &HouseBuilder{house: &House{}}
}

func (b *HouseBuilder) SetWalls(w int) *HouseBuilder {
	b.house.walls = w
	return b
}

func (b *HouseBuilder) SetDoors(d int) *HouseBuilder {
	b.house.doors = d
	return b
}

func (b *HouseBuilder) SetWindows(w int) *HouseBuilder {
	b.house.windows = w
	return b
}

func (b *HouseBuilder) Build() *House {
	return b.house
}

func main() {
	builder := NewHouseBuilder()
	house := builder.SetWalls(4).SetDoors(2).SetWindows(6).Build()

	fmt.Printf("House: %d walls, %d doors, %d windows\n", house.walls, house.doors, house.windows)
}
```
Особенности в Go: Используем цепочку вызовов для удобства.

### 5. Prototype (Прототип)
**Назначение:** Создаёт новые объекты путём копирования существующего экземпляра.
**Пример:** Клонирование пользователя.
```go
package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func (u *User) Clone() *User {
	return &User{Name: u.Name, Age: u.Age}
}

func main() {
	user1 := &User{Name: "Alice", Age: 25}
	user2 := user1.Clone()

	user2.Name = "Bob"
	fmt.Println(user1) // &{Alice 25}
	fmt.Println(user2) // &{Bob 25}
}
```
Особенности в Go: Простое копирование структуры.

### 6. Adapter (Адаптер)
**Назначение:** Позволяет объектам с несовместимыми интерфейсами работать вместе.
**Пример:** Адаптация старого логгера к новому интерфейсу.
```go
package main

import "fmt"

type NewLogger interface {
	Log(msg string)
}

type OldLogger struct{}

func (o *OldLogger) WriteLog(msg string) {
	fmt.Println("Old Log:", msg)
}

type LoggerAdapter struct {
	oldLogger *OldLogger
}

func (a *LoggerAdapter) Log(msg string) {
	a.oldLogger.WriteLog(msg)
}

func main() {
	oldLogger := &OldLogger{}
	adapter := &LoggerAdapter{oldLogger: oldLogger}

	adapter.Log("Test message") // Old Log: Test message
}
```

### 7. Decorator (Декоратор)
**Назначение:** Добавляет новое поведение объекту динамически.
**Пример:** Добавление логирования к функции.
```go
package main

import "fmt"

type Operation func(int) int

func Double(n int) int {
	return n * 2
}

func WithLogging(op Operation) Operation {
	return func(n int) int {
		fmt.Println("Input:", n)
		result := op(n)
		fmt.Println("Output:", result)
		return result
	}
}

func main() {
	doubleWithLog := WithLogging(Double)
	result := doubleWithLog(5)
	fmt.Println("Result:", result)
}
```
Вывод:
Input: 5
Output: 10
Result: 10

### 8. Observer (Наблюдатель)
**Назначение:** Определяет зависимость "один-ко-многим" между объектами.
**Пример:** Уведомление подписчиков о новом посте.
```go
package main

import "fmt"

type Subject struct {
	observers []Observer
}

type Observer interface {
	Update(message string)
}

type User struct {
	name string
}

func (u *User) Update(message string) {
	fmt.Printf("%s received: %s\n", u.name, message)
}

func (s *Subject) AddObserver(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *Subject) Notify(message string) {
	for _, o := range s.observers {
		o.Update(message)
	}
}

func main() {
	subject := &Subject{}
	user1 := &User{name: "Alice"}
	user2 := &User{name: "Bob"}

	subject.AddObserver(user1)
	subject.AddObserver(user2)
	subject.Notify("New post!")
}
```
Вывод:
Alice received: New post!
Bob received: New post!

### 9. Strategy (Стратегия)
**Назначение:** Определяет семейство алгоритмов и позволяет менять их на лету.
**Пример:** Разные способы оплаты.
```go
package main

import "fmt"

type PaymentStrategy interface {
	Pay(amount int) string
}

type CreditCard struct{}
func (c *CreditCard) Pay(amount int) string { return fmt.Sprintf("Paid %d with Credit Card", amount) }

type PayPal struct{}
func (p *PayPal) Pay(amount int) string { return fmt.Sprintf("Paid %d with PayPal", amount) }

type PaymentContext struct {
	strategy PaymentStrategy
}

func (p *PaymentContext) SetStrategy(s PaymentStrategy) {
	p.strategy = s
}

func (p *PaymentContext) ExecutePayment(amount int) string {
	return p.strategy.Pay(amount)
}

func main() {
	payment := &PaymentContext{}
	payment.SetStrategy(&CreditCard{})
	fmt.Println(payment.ExecutePayment(100))

	payment.SetStrategy(&PayPal{})
	fmt.Println(payment.ExecutePayment(200))
}
```

### 10. Command (Команда)
**Назначение:** Инкапсулирует запрос как объект.
**Пример:** Управление светом.
```go
package main

import "fmt"

type Command interface {
	Execute()
}

type Light struct{}

func (l *Light) On()  { fmt.Println("Light is ON") }
func (l *Light) Off() { fmt.Println("Light is OFF") }

type LightOnCommand struct {
	light *Light
}

func (c *LightOnCommand) Execute() {
	c.light.On()
}

type LightOffCommand struct {
	light *Light
}

func (c *LightOffCommand) Execute() {
	c.light.Off()
}

type Remote struct {
	command Command
}

func (r *Remote) SetCommand(c Command) {
	r.command = c
}

func (r *Remote) PressButton() {
	r.command.Execute()
}

func main() {
	light := &Light{}
	remote := &Remote{}

	remote.SetCommand(&LightOnCommand{light: light})
	remote.PressButton()

	remote.SetCommand(&LightOffCommand{light: light})
	remote.PressButton()
}
```
Вывод:
Light is ON
Light is OFF
