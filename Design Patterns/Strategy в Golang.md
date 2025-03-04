# Шаблон проектирования Strategy в Golang

## Введение

Шаблон проектирования **Strategy** (Стратегия) — это поведенческий шаблон, который определяет семейство алгоритмов, инкапсулирует каждый из них и делает их взаимозаменяемыми внутри контекста. В Go, где нет классов в традиционном смысле, Strategy реализуется через интерфейсы, структуры и функции, что делает его подходящим для динамического выбора алгоритма во время выполнения.

В этой лекции мы разберём:
- Что такое Strategy и где он применяется.
- Как реализовать Strategy в Go.
- Преимущества и недостатки шаблона.
- Примеры использования в реальных задачах.
- Рекомендации по применению в Go.

---

## 1. Что такое Strategy?

Strategy — это шаблон, который:
- Позволяет определять разные алгоритмы (стратегии) для выполнения одной и той же задачи.
- Инкапсулирует эти алгоритмы, делая их взаимозаменяемыми.
- Предоставляет контексту возможность динамически выбирать стратегию во время выполнения.
- Снижает связанность между клиентским кодом и конкретными реализациями алгоритмов.

### Примеры использования:
- Сортировка данных (разные алгоритмы: быстрая сортировка, пузырьковая, слияние).
- Форматирование данных (JSON, XML, YAML).
- Обработка платежей (кредитная карта, PayPal, наличные).
- Сжатие данных (Gzip, Zip, без сжатия).

---

## 2. Реализация Strategy в Go

В Go Strategy реализуется через интерфейсы и структуры, где стратегия представлена интерфейсом, а конкретные алгоритмы — структурами, реализующими этот интерфейс. Контекст содержит ссылку на стратегию и использует её для выполнения операций.

### 2.1. Базовая структура

Рассмотрим пример системы сортировки, где можно использовать разные стратегии сортировки (пузырьковая и быстрая сортировка).

#### Шаг 1: Определение интерфейса стратегии
```go
package strategy

// SortStrategy — интерфейс для стратегий сортировки
type SortStrategy interface {
    Sort(data []int) []int
}
```

#### Шаг 2: Конкретные стратегии
Создадим две реализации: `BubbleSort` (пузырьковая сортировка) и `QuickSort` (быстрая сортировка):

```go
// BubbleSort — реализация пузырьковой сортировки
type BubbleSort struct{}

func (b *BubbleSort) Sort(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    n := len(result)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    return result
}

// QuickSort — реализация быстрой сортировки
type QuickSort struct{}

func (q *QuickSort) Sort(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    quickSort(result, 0, len(result)-1)
    return result
}

func quickSort(arr []int, low, high int) {
    if low < high {
        pivot := partition(arr, low, high)
        quickSort(arr, low, pivot-1)
        quickSort(arr, pivot+1, high)
    }
}

func partition(arr []int, low, high int) int {
    pivot := arr[high]
    i := low - 1
    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}
```

#### Шаг 3: Контекст (Sorter)
Создадим структуру, которая использует стратегию для сортировки:

```go
// Sorter — контекст, использующий стратегию сортировки
type Sorter struct {
    strategy SortStrategy
}

// NewSorter — конструктор для создания Sorter
func NewSorter(strategy SortStrategy) *Sorter {
    return &Sorter{strategy: strategy}
}

// SetStrategy — изменение стратегии
func (s *Sorter) SetStrategy(strategy SortStrategy) {
    s.strategy = strategy
}

// Sort — выполнение сортировки с выбранной стратегией
func (s *Sorter) Sort(data []int) []int {
    return s.strategy.Sort(data)
}
```

#### Шаг 4: Использование
```go
package main

import (
    "fmt"
    "strategy"
)

func main() {
    data := []int{64, 34, 25, 12, 22, 11, 90}

    // Используем пузырьковую сортировку
    bubbleSorter := strategy.NewSorter(&strategy.BubbleSort{})
    sortedBubble := bubbleSorter.Sort(data)
    fmt.Println("Пузырьковая сортировка:", sortedBubble)

    // Меняем стратегию на быструю сортировку
    bubbleSorter.SetStrategy(&strategy.QuickSort{})
    sortedQuick := bubbleSorter.Sort(data)
    fmt.Println("Быстрая сортировка:", sortedQuick)
}
```

**Вывод:**
```
Пузырьковая сортировка: [11 12 22 25 34 64 90]
Быстрая сортировка: [11 12 22 25 34 64 90]
```

---

### 2.2. Расширенная реализация с конфигурацией

Добавим возможность передавать конфигурацию для выбора стратегии и поддержку дополнительных алгоритмов:

```go
package strategy

// SortConfig — конфигурация для сортировки
type SortConfig struct {
    Algorithm string
    Threshold int // Пример дополнительного параметра
}

// SortStrategy — интерфейс для стратегий сортировки
type SortStrategy interface {
    Sort(data []int, config SortConfig) []int
}

// BubbleSort — реализация пузырьковой сортировки
type BubbleSort struct{}

func (b *BubbleSort) Sort(data []int, config SortConfig) []int {
    result := make([]int, len(data))
    copy(result, data)
    n := len(result)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    return result
}

// QuickSort — реализация быстрой сортировки
type QuickSort struct{}

func (q *QuickSort) Sort(data []int, config SortConfig) []int {
    result := make([]int, len(data))
    copy(result, data)
    quickSort(result, 0, len(result)-1, config.Threshold)
    return result
}

func quickSort(arr []int, low, high, threshold int) {
    if low < high {
        if high-low < threshold { // Используем порог для оптимизации
            bubbleSort(arr, low, high)
        } else {
            pivot := partition(arr, low, high)
            quickSort(arr, low, pivot-1, threshold)
            quickSort(arr, pivot+1, high, threshold)
        }
    }
}

func bubbleSort(arr []int, low, high int) {
    for i := low; i < high; i++ {
        for j := low; j < high-i+low; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
            }
        }
    }
}

func partition(arr []int, low, high int) int {
    pivot := arr[high]
    i := low - 1
    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}

// Sorter — контекст, использующий стратегию сортировки
type Sorter struct {
    strategy SortStrategy
}

// NewSorter — конструктор для создания Sorter
func NewSorter(strategy SortStrategy) *Sorter {
    return &Sorter{strategy: strategy}
}

// SetStrategy — изменение стратегии
func (s *Sorter) SetStrategy(strategy SortStrategy) {
    s.strategy = strategy
}

// Sort — выполнение сортировки с выбранной стратегией и конфигурацией
func (s *Sorter) Sort(data []int, config SortConfig) []int {
    return s.strategy.Sort(data, config)
}

// CreateStrategy — фабричный метод для создания стратегии
func CreateStrategy(config SortConfig) SortStrategy {
    switch config.Algorithm {
    case "bubble":
        return &BubbleSort{}
    case "quick":
        return &QuickSort{}
    default:
        return nil
    }
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "strategy"
)

func main() {
    data := []int{64, 34, 25, 12, 22, 11, 90}

    // Используем пузырьковую сортировку
    configBubble := strategy.SortConfig{Algorithm: "bubble", Threshold: 5}
    sorter := strategy.NewSorter(strategy.CreateStrategy(configBubble))
    sortedBubble := sorter.Sort(data, configBubble)
    fmt.Println("Пузырьковая сортировка:", sortedBubble)

    // Меняем стратегию на быструю сортировку
    configQuick := strategy.SortConfig{Algorithm: "quick", Threshold: 5}
    sorter.SetStrategy(strategy.CreateStrategy(configQuick))
    sortedQuick := sorter.Sort(data, configQuick)
    fmt.Println("Быстрая сортировка:", sortedQuick)
}
```

**Вывод:**
```
Пузырьковая сортировка: [11 12 22 25 34 64 90]
Быстрая сортировка: [11 12 22 25 34 64 90]
```

---

## 3. Преимущества Strategy

- **Гибкость**: Позволяет динамически менять алгоритмы во время выполнения без изменения клиентского кода.
- **Инкапсуляция**: Скрывает реализацию алгоритмов, упрощая клиентский код.
- **Расширяемость**: Легко добавлять новые стратегии без изменения существующего кода.
- **Соответствие принципам SOLID**: Следует принципу открытости/закрытости (Open/Closed Principle).

---

## 4. Недостатки Strategy

- **Усложнение кода**: Множество стратегий может сделать код сложным для понимания и отладки.
- **Производительность**: Переключение между стратегиями может добавить накладные расходы.
- **Сложности с отладкой**: Трудно отслеживать, какая стратегия используется, если их много.
- **Избыточность**: Для простых случаев может быть избыточным, когда достаточно прямого вызова функции.

---

## 5. Примеры реального использования

### 5.1. Обработка платежей
Реализация разных стратегий оплаты (кредитная карта, PayPal):

```go
package payment

import "fmt"

// PaymentStrategy — интерфейс для стратегий оплаты
type PaymentStrategy interface {
    Pay(amount float64) error
}

// CreditCardStrategy — стратегия оплаты кредитной картой
type CreditCardStrategy struct {
    cardNumber string
}

func (c *CreditCardStrategy) Pay(amount float64) error {
    return fmt.Errorf("оплата %f с кредитной карты %s успешно проведена", amount, c.cardNumber)
}

// PayPalStrategy — стратегия оплаты через PayPal
type PayPalStrategy struct {
    email string
}

func (p *PayPalStrategy) Pay(amount float64) error {
    return fmt.Errorf("оплата %f через PayPal (%s) успешно проведена", amount, p.email)
}

// PaymentContext — контекст для обработки платежей
type PaymentContext struct {
    strategy PaymentStrategy
}

// NewPaymentContext — конструктор для создания контекста
func NewPaymentContext(strategy PaymentStrategy) *PaymentContext {
    return &PaymentContext{strategy: strategy}
}

// SetStrategy — изменение стратегии
func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}

// Pay — выполнение оплаты с выбранной стратегией
func (p *PaymentContext) Pay(amount float64) error {
    return p.strategy.Pay(amount)
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "payment"
)

func main() {
    // Оплата кредитной картой
    creditCard := &payment.CreditCardStrategy{cardNumber: "1234-5678-9012-3456"}
    context := payment.NewPaymentContext(creditCard)
    err := context.Pay(100.50)
    if err != nil {
        fmt.Println(err)
    }

    // Переключение на PayPal
    payPal := &payment.PayPalStrategy{email: "user@example.com"}
    context.SetStrategy(payPal)
    err = context.Pay(200.00)
    if err != nil {
        fmt.Println(err)
    }
}
```

**Вывод:**
```
оплата 100.5 с кредитной карты 1234-5678-9012-3456 успешно проведена
оплата 200 через PayPal (user@example.com) успешно проведена
```

---

### 5.2. Форматирование данных
Реализация разных стратегий форматирования (JSON, XML):

```go
package formatter

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
)

// FormatStrategy — интерфейс для стратегий форматирования
type FormatStrategy interface {
    Format(data interface{}) (string, error)
}

// JSONStrategy — стратегия форматирования в JSON
type JSONStrategy struct{}

func (j *JSONStrategy) Format(data interface{}) (string, error) {
    bytes, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// XMLStrategy — стратегия форматирования в XML
type XMLStrategy struct{}

func (x *XMLStrategy) Format(data interface{}) (string, error) {
    bytes, err := xml.MarshalIndent(data, "", "  ")
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// FormatterContext — контекст для форматирования данных
type FormatterContext struct {
    strategy FormatStrategy
}

// NewFormatterContext — конструктор для создания контекста
func NewFormatterContext(strategy FormatStrategy) *FormatterContext {
    return &FormatterContext{strategy: strategy}
}

// SetStrategy — изменение стратегии
func (f *FormatterContext) SetStrategy(strategy FormatStrategy) {
    f.strategy = strategy
}

// Format — выполнение форматирования с выбранной стратегией
func (f *FormatterContext) Format(data interface{}) (string, error) {
    return f.strategy.Format(data)
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "formatter"
)

func main() {
    data := struct {
        Name string
        Age  int
    }{Name: "Иван", Age: 30}

    // Форматирование в JSON
    jsonFormatter := formatter.NewFormatterContext(&formatter.JSONStrategy{})
    jsonResult, err := jsonFormatter.Format(data)
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
    } else {
        fmt.Println("JSON:")
        fmt.Println(jsonResult)
    }

    // Переключение на XML
    jsonFormatter.SetStrategy(&formatter.XMLStrategy{})
    xmlResult, err := jsonFormatter.Format(data)
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
    } else {
        fmt.Println("XML:")
        fmt.Println(xmlResult)
    }
}
```

**Вывод (примерный):**
```
JSON:
{
  "Name": "Иван",
  "Age": 30
}
XML:
<struct>
  <Name>Иван</Name>
  <Age>30</Age>
</struct>
```

---

## 6. Рекомендации по использованию Strategy в Go

1. **Используйте интерфейсы**: Определите интерфейс `Strategy`, чтобы обеспечить гибкость и расширяемость.
2. **Избегайте избыточности**: Для простых случаев прямой вызов функции может быть проще, чем создание стратегии.
3. **Тестирование**: Создавайте мок-объекты для стратегий, чтобы легко тестировать контекст.
4. **Производительность**: Оценивайте влияние переключения стратегий на производительность, особенно для критичных операций.
5. **Читаемость**: Сделайте код стратегий и контекста понятным, добавляя документацию и простые имена.

---

## 7. Преимущества и недостатки

### Преимущества:
- **Гибкость**: Легко менять алгоритмы во время выполнения без изменения клиентского кода.
- **Инкапсуляция**: Скрывает реализацию алгоритмов, упрощая клиентский код.
- **Расширяемость**: Следует принципу открытости/закрытости (Open/Closed Principle).
- **Тестируемость**: Легко заменять стратегии на тестовые реализации (mocks).

### Недостатки:
- **Усложнение кода**: Множество стратегий может сделать код сложным для понимания.
- **Производительность**: Переключение между стратегиями может добавить накладные расходы.
- **Сложности с отладкой**: Трудно отслеживать, какая стратегия используется, если их много.
- **Избыточность**: Может быть избыточным для простых случаев.

---

## 8. Заключение

Шаблон Strategy в Go — мощный инструмент для динамического выбора алгоритмов. В Go он особенно удобен благодаря поддержке интерфейсов и функций. Используйте Strategy, когда нужно обеспечить гибкость в выборе алгоритма без изменения существующего кода, но избегайте его для простых случаев, где достаточно прямого вызова функции.