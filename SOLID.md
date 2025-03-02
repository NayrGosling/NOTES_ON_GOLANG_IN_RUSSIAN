# Принципы SOLID

## Что такое SOLID?

SOLID — это аббревиатура, обозначающая пять принципов объектно-ориентированного программирования, предложенных Робертом Мартином. Эти принципы помогают делать код более понятным, гибким и поддерживаемым.

## Принципы SOLID

### S - Single Responsibility Principle (Принцип единственной ответственности)

Каждый класс или модуль должен иметь только одну причину для изменения, то есть выполнять только одну задачу.

**Пример:**

```go
// Плохая практика: класс выполняет несколько задач
type Report struct {
    Content string
}

func (r *Report) Generate() string {
    return "Report content: " + r.Content
}

func (r *Report) SaveToFile(filename string) {
    os.WriteFile(filename, []byte(r.Content), 0644)
}

// Хорошая практика: разделение ответственности
// Класс отвечает только за генерацию отчета
```

### O - Open/Closed Principle (Принцип открытости/закрытости)

Классы должны быть открыты для расширения (например, через наследование или интерфейсы), но закрыты для модификации (изменения исходного кода).

**Пример:**

```go
// Нарушение принципа: добавление нового типа скидки требует изменения кода
func CalculateDiscount(price float64, discountType string) float64 {
    if discountType == "fixed" {
        return price - 10
    } else if discountType == "percentage" {
        return price * 0.9
    }
    return price
}

func main() {
    price := 100.0
    fmt.Println("Fixed discount:", CalculateDiscount(price, "fixed"))
    fmt.Println("Percentage discount:", CalculateDiscount(price, "percentage"))
}
```
```go
// Следование принципу: использование интерфейсов
// Интерфейс для вычисления скидки
type DiscountStrategy interface {
    ApplyDiscount(price float64) float64
}

// Конкретная реализация для фиксированной скидки
type FixedDiscount struct {
    Amount float64
}

func (fd FixedDiscount) ApplyDiscount(price float64) float64 {
    return price - fd.Amount
}

// Конкретная реализация для процентной скидки
type PercentageDiscount struct {
    Percentage float64
}

func (pd PercentageDiscount) ApplyDiscount(price float64) float64 {
    return price * (1 - pd.Percentage/100)
}

// Функция для применения скидки, принимающая любой тип, реализующий интерфейс DiscountStrategy
func CalculateDiscount(price float64, strategy DiscountStrategy) float64 {
    return strategy.ApplyDiscount(price)
}

func main() {
    price := 100.0

    // Применение фиксированной скидки
    fixedDiscount := FixedDiscount{Amount: 10}
    fmt.Println("Fixed discount:", CalculateDiscount(price, fixedDiscount))

    // Применение процентной скидки
    percentageDiscount := PercentageDiscount{Percentage: 10}
    fmt.Println("Percentage discount:", CalculateDiscount(price, percentageDiscount))
}
```

### L - Liskov Substitution Principle (Принцип подстановки Барбары Лисков)

Объекты базового класса должны заменяться объектами производных классов без нарушения работы программы.

**Пример:**

```go
// Базовый интерфейс "Shape"
type Shape interface {
	Area() float64
}

// Класс "Rectangle" (Прямоугольник), который реализует интерфейс "Shape"
type Rectangle struct {
	Width, Height float64
}

func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Класс "Square" (Квадрат), который также реализует интерфейс "Shape"
type Square struct {
	Side float64
}

func (s *Square) Area() float64 {
	return s.Side * s.Side
}

// Функция для печати площади фигуры
func printArea(s Shape) {
	fmt.Printf("Area: %f\n", s.Area())
}

func main() {
	// Создаем прямоугольник и квадрат
	rect := &Rectangle{Width: 5, Height: 10}
	square := &Square{Side: 4}

	// Используем их в функции, принимающей интерфейс Shape
	printArea(rect)
	printArea(square)
}
```

### I - Interface Segregation Principle (Принцип разделения интерфейсов)

Клиенты не должны зависеть от интерфейсов, которые они не используют. Лучше иметь много узкоспециализированных интерфейсов, чем один универсальный.

**Пример:**

```go
// Плохая практика: общий интерфейс заставляет реализовывать ненужные методы
type Worker interface {
	Work() string
	Eat() string
}

type Engineer struct{}
type Chef struct{}

// Инженер должен только работать, но также обязан реализовать Eat, что ему не нужно
func (e Engineer) Work() string {
	return "Engineer is working"
}

func (e Engineer) Eat() string {
	return "Engineer is eating"
}

// Повар должен только есть, но также обязан реализовать Work, что ему не нужно
func (c Chef) Work() string {
	return "Chef is working"
}

func (c Chef) Eat() string {
	return "Chef is eating"
}

func main() {
	engineer := Engineer{}
	chef := Chef{}

	fmt.Println(engineer.Work()) // Работает
	fmt.Println(engineer.Eat())  // Работает, но ненужно для инженера
	fmt.Println(chef.Work())     // Работает
	fmt.Println(chef.Eat())      // Работает, но ненужно для повара
}
```
```go
// Хорошая практика: интерфейсы разделены по назначению
type Workable interface {
	Work() string
}

type Eatable interface {
	Eat() string
}

type Engineer struct{}
type Chef struct{}

// Инженер теперь реализует только интерфейс Workable
func (e Engineer) Work() string {
	return "Engineer is working"
}

// Повар теперь реализует только интерфейс Eatable
func (c Chef) Eat() string {
	return "Chef is eating"
}

func main() {
	engineer := Engineer{}
	chef := Chef{}

	// Инженер работает
	fmt.Println(engineer.Work())

	// Повар ест
	fmt.Println(chef.Eat())
}
```
#### Объяснение:
- В плохом примере интерфейс Worker включает оба метода: Work и Eat. Это приводит к тому, что клиенты, такие как Engineer и Chef, обязаны реализовывать методы, которые им не нужны. Инженеру нужно только работать, а повару — только есть.
- В хорошем примере интерфейсы разделены: Workable отвечает за работу, а Eatable — за еду. Таким образом, каждый тип реализует только тот интерфейс, который ему действительно нужен, что соответствует принципу разделения интерфейсов.

### D - Dependency Inversion Principle (Принцип инверсии зависимостей)

Модули высокого уровня не должны зависеть от модулей низкого уровня. Оба должны зависеть от абстракций (интерфейсов), а не от конкретных реализаций.

**Пример:**

```go
// Плохая практика: зависимость от конкретного класса
// Класс с высокой зависимостью от конкретной реализации
package main

import "fmt"

type MySQLDatabase struct {}

func (db MySQLDatabase) Connect() {
    fmt.Println("Подключение к MySQL")
}

type Service struct {
    db MySQLDatabase
}

func (s Service) GetData() {
    s.db.Connect()
    fmt.Println("Получение данных")
}

func main() {
    service := Service{db: MySQLDatabase{}}
    service.GetData()
}
```
```go
Следование принципу (использование абстракции):
go
Копировать
Редактировать
package main

import "fmt"

// Абстракция для базы данных
type Database interface {
    Connect()
}

// Конкретная реализация для MySQL
type MySQLDatabase struct {}

func (db MySQLDatabase) Connect() {
    fmt.Println("Подключение к MySQL")
}

// Конкретная реализация для PostgreSQL
type PostgreSQLDatabase struct {}

func (db PostgreSQLDatabase) Connect() {
    fmt.Println("Подключение к PostgreSQL")
}

// Класс сервис, теперь зависимый от абстракции
type Service struct {
    db Database // Зависимость от абстракции, а не конкретной реализации
}

func (s Service) GetData() {
    s.db.Connect()
    fmt.Println("Получение данных")
}

func main() {
    // Подключение с использованием MySQL
    mySQLService := Service{db: MySQLDatabase{}}
    mySQLService.GetData()

    // Подключение с использованием PostgreSQL
    postgreSQLService := Service{db: PostgreSQLDatabase{}}
    postgreSQLService.GetData()
}
```
#### Объяснение:
- В плохом примере класс Service зависит напрямую от конкретной реализации базы данных MySQLDatabase, что нарушает принцип инверсии зависимостей.
- В хорошем примере создается абстракция Database, которая определяет метод Connect(). Классы MySQLDatabase и PostgreSQLDatabase реализуют этот интерфейс, а класс Service теперь зависит от абстракции, а не от конкретной реализации. Это позволяет легко заменять реализацию базы данных без изменения кода в классе Service.

## Заключение

Следование принципам SOLID делает код более структурированным, удобным для расширения и сопровождаемым. Эти принципы широко используются в разработке ПО для построения гибких и масштабируемых систем.

