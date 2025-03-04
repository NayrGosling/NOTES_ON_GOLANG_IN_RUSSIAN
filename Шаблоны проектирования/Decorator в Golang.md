# Шаблон проектирования Decorator в Golang

## Введение

Шаблон проектирования **Decorator** (Декоратор) — это структурный шаблон, который позволяет динамически добавлять новые поведения или функциональность к существующему объекту, не изменяя его код. В Go, где нет классов в традиционном смысле, Decorator реализуется через интерфейсы, структуры и композицию, что делает его подходящим для добавления функций "обёрткой" (wrapper) вокруг объектов.

В этой лекции мы разберём:
- Что такое Decorator и где он применяется.
- Как реализовать Decorator в Go.
- Преимущества и недостатки шаблона.
- Примеры использования в реальных задачах.
- Рекомендации по применению в Go.

---

## 1. Что такое Decorator?

Decorator — это шаблон, который:
- Позволяет добавлять новые обязанности (поведения) к объекту динамически, в процессе выполнения программы.
- Использует композицию вместо наследования, обёртывая объект в другой объект, который расширяет его функциональность.
- Сохраняет открытость для добавления новых декораторов без изменения существующего кода.

### Примеры использования:
- Добавление логов к методам (например, логирование вызовов функций).
- Кэширование результатов операций.
- Добавление сжатия или шифрования к потокам данных (например, HTTP-запросам).
- Обогащение функциональности UI-компонентов.

---

## 2. Реализация Decorator в Go

В Go Decorator реализуется через интерфейсы и структуры, где декоратор "обёртывает" базовый объект, реализуя тот же интерфейс и добавляя новое поведение. Go не использует классы или наследование, поэтому мы опираемся на композицию.

### 2.1. Базовая структура

Рассмотрим пример декоратора для объекта, представляющего кофе (Coffee). Мы добавим возможность добавления добавок (например, молока, сахара), которые увеличивают стоимость и меняют описание.

#### Шаг 1: Определение интерфейса
```go
package decorator

// Beverage — интерфейс для напитков
type Beverage interface {
    Cost() float64
    Description() string
}
```

#### Шаг 2: Базовая реализация (SimpleCoffee)
```go
// SimpleCoffee — базовый объект (простой кофе)
type SimpleCoffee struct{}

func (c *SimpleCoffee) Cost() float64 {
    return 2.0
}

func (c *SimpleCoffee) Description() string {
    return "Простой кофе"
}
```

#### Шаг 3: Декоратор (Decorator)
Создадим декораторы для добавления молока и сахара:

```go
// BeverageDecorator — базовая структура для декораторов
type BeverageDecorator struct {
    beverage Beverage
}

func (d *BeverageDecorator) Cost() float64 {
    return d.beverage.Cost()
}

func (d *BeverageDecorator) Description() string {
    return d.beverage.Description()
}

// MilkDecorator — декоратор для добавления молока
type MilkDecorator struct {
    BeverageDecorator
}

func NewMilkDecorator(beverage Beverage) *MilkDecorator {
    return &MilkDecorator{
        BeverageDecorator{beverage},
    }
}

func (m *MilkDecorator) Cost() float64 {
    return m.BeverageDecorator.Cost() + 0.5
}

func (m *MilkDecorator) Description() string {
    return m.BeverageDecorator.Description() + ", с молоком"
}

// SugarDecorator — декоратор для добавления сахара
type SugarDecorator struct {
    BeverageDecorator
}

func NewSugarDecorator(beverage Beverage) *SugarDecorator {
    return &SugarDecorator{
        BeverageDecorator{beverage},
    }
}

func (s *SugarDecorator) Cost() float64 {
    return s.BeverageDecorator.Cost() + 0.2
}

func (s *SugarDecorator) Description() string {
    return s.BeverageDecorator.Description() + ", с сахаром"
}
```

#### Шаг 4: Использование
```go
package main

import (
    "fmt"
    "decorator"
)

func main() {
    // Создаём простой кофе
    coffee := &decorator.SimpleCoffee{}
    fmt.Printf("Простой кофе: %s, Цена: $%.2f\n", coffee.Description(), coffee.Cost())

    // Добавляем молоко
    coffeeWithMilk := decorator.NewMilkDecorator(coffee)
    fmt.Printf("Кофе с молоком: %s, Цена: $%.2f\n", coffeeWithMilk.Description(), coffeeWithMilk.Cost())

    // Добавляем молоко и сахар
    coffeeWithMilkAndSugar := decorator.NewSugarDecorator(coffeeWithMilk)
    fmt.Printf("Кофе с молоком и сахаром: %s, Цена: $%.2f\n", coffeeWithMilkAndSugar.Description(), coffeeWithMilkAndSugar.Cost())
}
```

**Вывод:**
```
Простой кофе: Простой кофе, Цена: $2.00
Кофе с молоком: Простой кофе, с молоком, Цена: $2.50
Кофе с молоком и сахаром: Простой кофе, с молоком, с сахаром, Цена: $2.70
```

---

### 2.2. Расширенная реализация с логгированием

Добавим декоратор для логирования вызовов методов:

```go
package decorator

import (
    "fmt"
    "time"
)

// LoggedDecorator — декоратор для логирования вызовов
type LoggedDecorator struct {
    BeverageDecorator
}

func NewLoggedDecorator(beverage Beverage) *LoggedDecorator {
    return &LoggedDecorator{
        BeverageDecorator{beverage},
    }
}

func (l *LoggedDecorator) Cost() float64 {
    start := time.Now()
    cost := l.BeverageDecorator.Cost()
    duration := time.Since(start)
    fmt.Printf("Логирование: Вызов Cost занял %v, результат: $%.2f\n", duration, cost)
    return cost
}

func (l *LoggedDecorator) Description() string {
    start := time.Now()
    desc := l.BeverageDecorator.Description()
    duration := time.Since(start)
    fmt.Printf("Логирование: Вызов Description занял %v, результат: %s\n", duration, desc)
    return desc
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "decorator"
)

func main() {
    coffee := &decorator.SimpleCoffee{}
    loggedCoffee := decorator.NewLoggedDecorator(coffee)

    fmt.Printf("Описание: %s\n", loggedCoffee.Description())
    fmt.Printf("Цена: $%.2f\n", loggedCoffee.Cost())
}
```

**Вывод (примерный, время может варьироваться):**
```
Логирование: Вызов Description занял 1ns, результат: Простой кофе
Описание: Простой кофе
Логирование: Вызов Cost занял 1ns, результат: $2.00
Цена: $2.00
```

---

## 3. Преимущества Decorator

- **Гибкость**: Позволяет динамически добавлять новые поведения без изменения существующих объектов.
- **Инкапсуляция**: Скрывает добавляемую функциональность, делая код чище.
- **Расширяемость**: Легко добавлять новые декораторы без изменения базового кода.
- **Соответствие принципам SOLID**: Следует принципу открытости/закрытости (Open/Closed Principle).

---

## 4. Недостатки Decorator

- **Усложнение кода**: Множество декораторов может сделать код сложным для понимания и отладки.
- **Производительность**: Дополнительные обёртки могут немного замедлить выполнение (хотя в большинстве случаев это незначительно).
- **Сложности с отладкой**: Трудно отслеживать, какие декораторы применены, если их много.
- **Избыточность**: Для простых случаев может быть избыточным, когда достаточно прямого расширения объекта.

---

## 5. Примеры реального использования

### 5.1. Декоратор для HTTP-клиента с кэшированием
Добавим кэширование к HTTP-запросам:

```go
package httpdecorator

import (
    "net/http"
    "time"
)

// Client — интерфейс для HTTP-клиента
type Client interface {
    Get(url string) (*http.Response, error)
}

// BaseClient — базовый клиент
type BaseClient struct{}

func (c *BaseClient) Get(url string) (*http.Response, error) {
    return http.Get(url)
}

// CacheDecorator — декоратор для кэширования запросов
type CacheDecorator struct {
    client Client
    cache  map[string]*http.Response
}

func NewCacheDecorator(client Client) *CacheDecorator {
    return &CacheDecorator{
        client: client,
        cache:  make(map[string]*http.Response),
    }
}

func (c *CacheDecorator) Get(url string) (*http.Response, error) {
    if resp, ok := c.cache[url]; ok {
        return resp, nil // Возвращаем кэшированный ответ
    }
    resp, err := c.client.Get(url)
    if err == nil {
        c.cache[url] = resp // Кэшируем ответ
    }
    return resp, err
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "httpdecorator"
    "net/http"
)

func main() {
    baseClient := &httpdecorator.BaseClient{}
    cachedClient := httpdecorator.NewCacheDecorator(baseClient)

    // Первый запрос — реальный
    resp1, err1 := cachedClient.Get("https://api.example.com/data")
    if err1 != nil {
        fmt.Printf("Ошибка: %v\n", err1)
    } else {
        fmt.Println("Первый запрос выполнен")
    }

    // Второй запрос — из кэша
    resp2, err2 := cachedClient.Get("https://api.example.com/data")
    if err2 != nil {
        fmt.Printf("Ошибка: %v\n", err2)
    } else {
        fmt.Println("Второй запрос выполнен (из кэша)")
    }
}
```

---

### 5.2. Декоратор для логирования
Добавим логирование к функциям:

```go
package logdecorator

import (
    "fmt"
    "time"
)

// Service — интерфейс для сервиса
type Service interface {
    Process(data string) string
}

// BaseService — базовый сервис
type BaseService struct{}

func (s *BaseService) Process(data string) string {
    return fmt.Sprintf("Обработано: %s", data)
}

// LoggedDecorator — декоратор для логирования
type LoggedDecorator struct {
    service Service
}

func NewLoggedDecorator(service Service) *LoggedDecorator {
    return &LoggedDecorator{service: service}
}

func (l *LoggedDecorator) Process(data string) string {
    start := time.Now()
    result := l.service.Process(data)
    duration := time.Since(start)
    fmt.Printf("Логирование: Обработка '%s' заняла %v, результат: %s\n", data, duration, result)
    return result
}
```

#### Использование:
```go
package main

import (
    "logdecorator"
)

func main() {
    baseService := &logdecorator.BaseService{}
    loggedService := logdecorator.NewLoggedDecorator(baseService)

    result := loggedService.Process("тестовые данные")
    fmt.Println("Итоговый результат:", result)
}
```

**Вывод (примерный, время может варьироваться):**
```
Логирование: Обработка 'тестовые данные' заняла 1ns, результат: Обработано: тестовые данные
Итоговый результат: Обработано: тестовые данные
```

---

## 6. Рекомендации по использованию Decorator в Go

1. **Используйте интерфейсы**: Определите интерфейс для декорируемых объектов, чтобы обеспечить гибкость и расширяемость.
2. **Избегайте избыточности**: Для простых случаев прямое расширение объекта может быть проще, чем создание декораторов.
3. **Тестирование**: Создавайте мок-объекты для декораторов, чтобы легко тестировать клиентский код.
4. **Производительность**: Оценивайте влияние декораторов на производительность, особенно если их много.
5. **Читаемость**: Сделайте код декораторов понятным, добавляя документацию и простые имена.

---

## 7. Преимущества и недостатки

### Преимущества:
- **Гибкость**: Легко добавлять новые поведения без изменения существующих объектов.
- **Инкапсуляция**: Скрывает добавляемую функциональность, упрощая клиентский код.
- **Расширяемость**: Следует принципу открытости/закрытости (Open/Closed Principle).
- **Композиция**: Использует композицию вместо наследования, что делает код более гибким.

### Недостатки:
- **Усложнение кода**: Множество декораторов может сделать код сложным для понимания.
- **Производительность**: Дополнительные обёртки могут замедлить выполнение.
- **Сложности с отладкой**: Трудно отслеживать, какие декораторы применены, если их много.
- **Избыточность**: Может быть избыточным для простых случаев.

---

## 8. Заключение

Шаблон Decorator в Go — мощный инструмент для динамического добавления функциональности к объектам. В Go он особенно удобен благодаря поддержке интерфейсов и композиции. Используйте Decorator, когда нужно расширять поведение объектов без изменения их кода, но избегайте его для простых случаев, где достаточно прямого изменения объекта или функций.