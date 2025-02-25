1. Singleton (Одиночка)
Описание: Гарантирует, что у класса существует только один экземпляр и предоставляет глобальную точку доступа к этому экземпляру. 
В Go одиночка часто реализуется через пакеты, так как Go не имеет встроенного механизма для классов.

Пример в Go:
Предположим, мы хотим создать логгер, который будет использоваться в приложении только один раз.

package logger
import "sync"

type Logger struct {
    messages []string
}

var instance *Logger
var once sync.Once

func GetInstance() *Logger {
    once.Do(func() {
        instance = &Logger{messages: make([]string, 0)}
    })
    return instance
}

func (l *Logger) Log(message string) {
    l.messages = append(l.messages, message)
}

func (l *Logger) GetMessages() []string {
    return l.messages
}

Использование:

package main
import (
    "fmt"
    "logger"
)

func main() {
    logger1 := logger.GetInstance()
    logger1.Log("Первый лог")

    logger2 := logger.GetInstance()
    logger2.Log("Второй лог")

    fmt.Println(logger1.GetMessages()) // ["Первый лог", "Второй лог"]
    fmt.Println(logger2.GetMessages()) // ["Первый лог", "Второй лог"]
    fmt.Println(logger1 == logger2)    // true (один и тот же экземпляр)
}
Особенности в Go: Вместо классов и приватных конструкторов Go использует пакеты и sync.Once для гарантированной инициализации единственного экземпляра. Это позволяет избежать лишней сложности.
______________________________________________________________________________________________
2. Factory (Фабрика)
Описание: Создаёт объекты без явного указания их конкретного класса. В Go фабричные функции часто используются для создания сложных структур или интерфейсов.

Пример в Go:
Предположим, мы хотим создавать разные типы транспортных средств (например, автомобили, самолёты).

package factory
import "fmt"

type Vehicle interface {
    Drive() string
}

type Car struct{}

func (c *Car) Drive() string {
    return "Еду на машине!"
}

type Airplane struct{}

func (a *Airplane) Drive() string {
    return "Лечу на самолёте!"
}

func CreateVehicle(vehicleType string) Vehicle {
    switch vehicleType {
    case "car":
        return &Car{}
    case "airplane":
        return &Airplane{}
    default:
        return nil
    }
}
Использование:

package main
import (
    "fmt"
    "factory"
)

func main() {
    car := factory.CreateVehicle("car")
    airplane := factory.CreateVehicle("airplane")

    fmt.Println(car.Drive())     // "Еду на машине!"
    fmt.Println(airplane.Drive()) // "Лечу на самолёте!"
}
Особенности в Go: Go не использует классы, поэтому фабрика реализуется через функции, возвращающие интерфейсы или структуры. Это делает код гибким и легко расширяемым.
______________________________________________________________________________________________
3. Observer (Наблюдатель)
Описание: Позволяет объекту уведомлять других объектов о изменениях состояния. В Go это часто реализуется через каналы (channels) или интерфейсы.

Пример в Go:
Предположим, у нас есть новостная система, где подписчики получают уведомления.

package observer
import "fmt"

type Subscriber interface {
    Notify(message string)
}

type NewsAgency struct {
    subscribers []Subscriber
}

func (n *NewsAgency) Register(subscriber Subscriber) {
    n.subscribers = append(n.subscribers, subscriber)
}

func (n *NewsAgency) Broadcast(message string) {
    for _, subscriber := range n.subscribers {
        subscriber.Notify(message)
    }
}

type User struct {
    name string
}

func (u *User) Notify(message string) {
    fmt.Printf("Пользователь %s получил новость: %s\n", u.name, message)
}
Использование:

package main
import ("observer")

func main() {
    agency := &observer.NewAgency{}
    
    user1 := &observer.User{name: "Алексей"}
    user2 := &observer.User{name: "Мария"}

    agency.Register(user1)
    agency.Register(user2)

    agency.Broadcast("Новая новость: Go 2.0 вышел!")
}
Вывод:
Пользователь Алексей получил новость: Новая новость: Go 2.0 вышел!
Пользователь Мария получил новость: Новая новость: Go 2.0 вышел!
Особенности в Go: Go использует интерфейсы и композицию вместо сложных наследований. Каналы также могут быть альтернативой для передачи уведомлений в реальном времени.
______________________________________________________________________________________________
4. Decorator (Декоратор)
Описание: Позволяет динамически добавлять ответственность к объекту, обёртывая его в другой объект. В Go это часто реализуется через композицию структур.

Пример в Go:
Предположим, мы хотим добавить функциональность к кофе (например, молоко, сахар).

package decorator
import "fmt"

type Beverage interface {
    Cost() float64
    Description() string
}

type BasicCoffee struct{}

func (b *BasicCoffee) Cost() float64 {
    return 2.0
}

func (b *BasicCoffee) Description() string {
    return "Простой кофе"
}

type MilkDecorator struct {
    beverage Beverage
}

func (m *MilkDecorator) Cost() float64 {
    return m.beverage.Cost() + 0.5
}

func (m *MilkDecorator) Description() string {
    return m.beverage.Description() + ", с молоком"
}

type SugarDecorator struct {
    beverage Beverage
}

func (s *SugarDecorator) Cost() float64 {
    return s.beverage.Cost() + 0.3
}

func (s *SugarDecorator) Description() string {
    return s.beverage.Description() + ", с сахаром"
}

Использование:

package main
import (
    "fmt"
    "decorator"
)

func main() {
    coffee := &decorator.BasicCoffee{}
    coffeeWithMilk := &decorator.MilkDecorator{beverage: coffee}
    coffeeWithMilkAndSugar := &decorator.SugarDecorator{beverage: coffeeWithMilk}

    fmt.Printf("Описание: %s, Цена: %.2f\n", coffeeWithMilkAndSugar.Description(), coffeeWithMilkAndSugar.Cost())
}
Вывод:
Описание: Простой кофе, с молоком, с сахаром, Цена: 2.80
Особенности в Go: Go не использует наследование для декораторов, а полагается на композицию через структуры, что делает код более гибким и читаемым.
______________________________________________________________________________________________
5. Strategy (Стратегия)
Описание: Позволяет определять семейство алгоритмов, инкапсулировать каждый из них и делать их взаимозаменяемыми. Клиентский код может выбирать алгоритм во время выполнения без изменения структуры программы.

Пример в Go:
Предположим, у нас есть система оплаты, где мы можем использовать разные стратегии оплаты (наличные, кредитная карта, PayPal).

package strategy
import "fmt"

// Интерфейс стратегии
type PaymentStrategy interface {
    Pay(amount float64) string
}

// Конкретная стратегия: Наличные
type CashPayment struct{}

func (c *CashPayment) Pay(amount float64) string {
    return fmt.Sprintf("Оплата наличными: %v рублей", amount)
}

// Конкретная стратегия: Кредитная карта
type CreditCardPayment struct{}

func (c *CreditCardPayment) Pay(amount float64) string {
    return fmt.Sprintf("Оплата кредитной картой: %v рублей", amount)
}

// Контекст, использующий стратегию
type ShoppingCart struct {
    paymentStrategy PaymentStrategy
}

func (s *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
    s.paymentStrategy = strategy
}

func (s *ShoppingCart) Checkout(amount float64) string {
    if s.paymentStrategy == nil {
        return "Стратегия оплаты не выбрана"
    }
    return s.paymentStrategy.Pay(amount)
}
Использование:

package main
import ("fmt", "strategy")

func main() {
    cart := &strategy.ShoppingCart{}

    // Оплата наличными
    cart.SetPaymentStrategy(&strategy.CashPayment{})
    fmt.Println(cart.Checkout(100.0)) // "Оплата наличными: 100 рублей"

    // Оплата кредитной картой
    cart.SetPaymentStrategy(&strategy.CreditCardPayment{})
    fmt.Println(cart.Checkout(200.0)) // "Оплата кредитной картой: 200 рублей"
}
Особенности в Go: Strategy реализуется через интерфейс (PaymentStrategy) и структуры, которые реализуют этот интерфейс. Это позволяет легко добавлять новые стратегии без изменения существующего кода.
______________________________________________________________________________________________
6. Adapter (Адаптер)
Описание: Позволяет объектам с несовместимыми интерфейсами работать вместе, адаптируя интерфейс одного объекта под другой.

Пример в Go:
Предположим, у нас есть старый сервис уведомлений с несовместимым интерфейсом, и мы хотим адаптировать его под новый интерфейс уведомлений.

package adapter
import "fmt"

// Целевой интерфейс (новый)
type Notification interface {
    Send(message string) string
}

// Адаптируемый класс (старый)
type OldNotificationService struct{}

func (o *OldNotificationService) Notify(message string) string {
    return fmt.Sprintf("Старый сервис: %s", message)
}

// Адаптер
type NotificationAdapter struct {
    oldService *OldNotificationService
}

func (a *NotificationAdapter) Send(message string) string {
    return a.oldService.Notify(message)
}
Использование:

package main
import ("fmt", "adapter")

func main() {
    oldService := &adapter.OldNotificationService{}
    adapter := &adapter.NotificationAdapter{oldService: oldService}

    // Используем адаптер как новый интерфейс Notification
    notifier := adapter.(adapter.Notification)
    fmt.Println(notifier.Send("Привет, мир!")) // "Старый сервис: Привет, мир!"
}
Особенности в Go: Адаптер реализуется через композицию (встраивание старого сервиса в новую структуру) и реализацию нового интерфейса. Go не требует сложных наследований, что упрощает создание адаптеров.
______________________________________________________________________________________________
7. Proxy (Прокси)
Описание: Предоставляет замену реальному объекту, чтобы контролировать доступ к нему, добавлять дополнительную логику (например, кэширование, ленивую загрузку, безопасность).

Пример в Go:
Предположим, у нас есть сервис изображений, и мы хотим добавить прокси для кэширования запросов.

package proxy
import "fmt"

// Реальный объект
type Image struct {
    name string
}

func (i *Image) Display() string {
    return fmt.Sprintf("Отображение изображения: %s", i.name)
}

// Прокси
type ImageProxy struct {
    realImage *Image
    name      string
}

func (p *ImageProxy) Display() string {
    if p.realImage == nil {
        p.realImage = &Image{name: p.name}
        fmt.Println("Загрузка изображения...")
    }
    return p.realImage.Display()
}
Использование:

package main
import ("fmt","proxy")

func main() {
    // Создаём прокси для изображения
    proxyImage := &proxy.ImageProxy{name: "photo.jpg"}

    // Первый вызов — изображение загружается
    fmt.Println(proxyImage.Display()) // "Загрузка изображения..." + "Отображение изображения: photo.jpg"

    // Повторный вызов — изображение уже загружено
    fmt.Println(proxyImage.Display()) // "Отображение изображения: photo.jpg"
}
Особенности в Go: Прокси реализуется через композицию, где прокси-объект содержит реальный объект и управляет доступом к нему. 
Go использует интерфейсы для обеспечения совместимости, но в простых случаях, как в примере, это может быть и структура с прямым вызовом методов.
