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

### 1. **Singleton (Одиночка)**

Назначение: Гарантирует, что у класса есть только один экземпляр, и предоставляет глобальную точку доступа к нему.

Пример: Логгер, который должен быть одним для всей системы.

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
Особенности в Go: Используем sync.Once для потокобезоп


### 2. **Factory Method (Фабричный метод)**

Назначение:
Определяет интерфейс для создания объекта, но позволяет подклассам решать, какой класс инстанцировать.

Пример:
Создание разных типов транспорта.

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
Особенности в Go:
Используем функцию вместо метода, так как в Go нет наследования.

### 3. **Abstract Factory (Абстрактная фабрика)**

Назначение:
Предоставляет интерфейс для создания семейств связанных объектов без указания их конкретных классов.

Пример:
Фабрика для создания UI-компонентов (кнопки и чекбоксы).

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
### 4. **Builder (Строитель)**

Назначение:
Отделяет конструирование сложного объекта от его представления.

Пример:
Постройка дома с разными конфигурациями.

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
Особенности в Go:
Используем цепочку вызовов для удобства.
