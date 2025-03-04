# Шаблон проектирования Builder в Golang

## Введение

Шаблон проектирования **Builder** (Строитель) — это порождающий шаблон, который позволяет создавать сложные объекты пошагово, разделяя процесс создания объекта от его представления. В Go, где нет классов в традиционном смысле, Builder реализуется через структуры, методы и интерфейсы, что делает его удобным для создания объектов с множеством опций или конфигураций.

В этой лекции мы разберём:
- Что такое Builder и где он применяется.
- Как реализовать Builder в Go.
- Преимущества и недостатки шаблона.
- Примеры использования в реальных задачах.
- Рекомендации по применению в Go.

---

## 1. Что такое Builder?

Builder — это шаблон, который:
- Разделяет создание сложного объекта на отдельные шаги.
- Позволяет создавать объекты с различными конфигурациями, не создавая множество конструкторов или перегруженных методов.
- Предоставляет чистый и читаемый способ сборки объектов с опциональными полями.

### Примеры использования:
- Создание сложных структур данных, таких как HTTP-запросы, конфигурации или объекты с множеством полей (например, пользователь с именем, email, ролями и т.д.).
- Построение объектов с поэтапной настройкой, например, запросов к API или базы данных.
- Генерация сложных документов (PDF, HTML) с разными настройками.

---

## 2. Реализация Builder в Go

В Go Builder реализуется через структуру (Builder) с методами, которые поэтапно настраивают объект, и конечный метод, возвращающий готовый объект. Go не использует классы или наследование, поэтому мы опираемся на композицию и методы.

### 2.1. Базовая структура

Рассмотрим пример построения объекта `User` с различными полями, которые могут быть опциональными.

#### Шаг 1: Определение целевой структуры
```go
package builder

import "time"

// User — конечный объект, который мы строим
type User struct {
    ID        int
    Name      string
    Email     string
    Age       int
    CreatedAt time.Time
}
```

#### Шаг 2: Создание строителя (Builder)
```go
// UserBuilder — структура для пошагового построения User
type UserBuilder struct {
    user User
}

// NewUserBuilder — конструктор для создания строителя
func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        user: User{
            CreatedAt: time.Now(), // Устанавливаем текущее время по умолчанию
        },
    }
}

// SetName — метод для установки имени
func (b *UserBuilder) SetName(name string) *UserBuilder {
    b.user.Name = name
    return b
}

// SetEmail — метод для установки email
func (b *UserBuilder) SetEmail(email string) *UserBuilder {
    b.user.Email = email
    return b
}

// SetAge — метод для установки возраста
func (b *UserBuilder) SetAge(age int) *UserBuilder {
    b.user.Age = age
    return b
}

// SetID — метод для установки ID
func (b *UserBuilder) SetID(id int) *UserBuilder {
    b.user.ID = id
    return b
}

// Build — конечный метод, возвращающий готовый объект
func (b *UserBuilder) Build() User {
    return b.user
}
```

#### Шаг 3: Использование
```go
package main

import (
    "fmt"
    "builder"
    "time"
)

func main() {
    // Создаём пользователя с минимальными данными
    user1 := builder.NewUserBuilder().
        SetName("Иван Иванов").
        SetEmail("ivan@example.com").
        Build()

    fmt.Printf("Пользователь 1: %+v\n", user1)

    // Создаём пользователя с дополнительными данными
    user2 := builder.NewUserBuilder().
        SetName("Мария Петрова").
        SetEmail("maria@example.com").
        SetAge(30).
        SetID(1).
        Build()

    fmt.Printf("Пользователь 2: %+v\n", user2)
}
```

**Вывод (примерный, время будет текущим):**
```
Пользователь 1: {ID:0 Name:Иван Иванов Email:ivan@example.com Age:0 CreatedAt:2025-03-03 12:00:00 +0000 UTC}
Пользователь 2: {ID:1 Name:Мария Петрова Email:maria@example.com Age:30 CreatedAt:2025-03-03 12:00:00 +0000 UTC}
```

---

### 2.2. Расширенная реализация с валидацией

Добавим валидацию и более сложную логику в Builder для проверки данных:

```go
package builder

import (
    "errors"
    "time"
)

// User — конечный объект, который мы строим
type User struct {
    ID        int
    Name      string
    Email     string
    Age       int
    CreatedAt time.Time
}

// UserBuilder — структура для пошагового построения User
type UserBuilder struct {
    user User
    err  error
}

// NewUserBuilder — конструктор для создания строителя
func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        user: User{
            CreatedAt: time.Now(),
        },
    }
}

// SetName — метод для установки имени с валидацией
func (b *UserBuilder) SetName(name string) *UserBuilder {
    if len(name) < 2 {
        b.err = errors.New("имя должно быть длиной не менее 2 символов")
        return b
    }
    b.user.Name = name
    return b
}

// SetEmail — метод для установки email с валидацией
func (b *UserBuilder) SetEmail(email string) *UserBuilder {
    if !isValidEmail(email) { // Простая проверка email
        b.err = errors.New("некорректный email")
        return b
    }
    b.user.Email = email
    return b
}

// SetAge — метод для установки возраста с валидацией
func (b *UserBuilder) SetAge(age int) *UserBuilder {
    if age < 0 || age > 150 {
        b.err = errors.New("возраст должен быть в диапазоне 0-150")
        return b
    }
    b.user.Age = age
    return b
}

// SetID — метод для установки ID
func (b *UserBuilder) SetID(id int) *UserBuilder {
    if id < 0 {
        b.err = errors.New("ID не может быть отрицательным")
        return b
    }
    b.user.ID = id
    return b
}

// Build — конечный метод, возвращающий готовый объект или ошибку
func (b *UserBuilder) Build() (User, error) {
    if b.err != nil {
        return User{}, b.err
    }
    return b.user, nil
}

// isValidEmail — простая функция для проверки email (для примера)
func isValidEmail(email string) bool {
    return len(email) > 0 && strings.Contains(email, "@")
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "builder"
)

func main() {
    // Успешное создание пользователя
    user1, err := builder.NewUserBuilder().
        SetName("Иван Иванов").
        SetEmail("ivan@example.com").
        SetAge(25).
        SetID(1).
        Build()

    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
    } else {
        fmt.Printf("Пользователь 1: %+v\n", user1)
    }

    // Попытка создать пользователя с ошибкой
    user2, err := builder.NewUserBuilder().
        SetName("A"). // Слишком короткое имя
        SetEmail("invalid-email"). // Некорректный email
        Build()

    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
    } else {
        fmt.Printf("Пользователь 2: %+v\n", user2)
    }
}
```

**Вывод (примерный):**
```
Пользователь 1: {ID:1 Name:Иван Иванов Email:ivan@example.com Age:25 CreatedAt:2025-03-03 12:00:00 +0000 UTC}
Ошибка: имя должно быть длиной не менее 2 символов
```

---

## 3. Преимущества Builder

- **Читаемость**: Позволяет создавать объекты с чёткой цепочкой вызовов, улучшая читаемость кода.
- **Гибкость**: Поддерживает опциональные поля и сложные конфигурации без необходимости создания множества конструкторов.
- **Валидация**: Легко добавлять проверки данных на каждом этапе построения.
- **Расширяемость**: Можно добавлять новые методы в Builder без изменения существующего кода.

---

## 4. Недостатки Builder

- **Усложнение кода**: Для простых объектов Builder может быть избыточным, увеличивая сложность.
- **Избыточность**: Требует создания дополнительной структуры (Builder), что может быть неоправданным для небольших объектов.
- **Сложности с тестированием**: Если Builder содержит сложную логику, тесты могут стать сложнее.

---

## 5. Примеры реального использования

### 5.1. Построение HTTP-запроса
Создание сложного HTTP-запроса с различными параметрами:

```go
package httpbuilder

import (
    "net/http"
    "time"
)

// Request — конечный объект HTTP-запроса
type Request struct {
    Method     string
    URL        string
    Headers    map[string]string
    Timeout    time.Duration
    Body       string
}

// RequestBuilder — структура для пошагового построения Request
type RequestBuilder struct {
    request Request
}

// NewRequestBuilder — конструктор для создания строителя
func NewRequestBuilder() *RequestBuilder {
    return &RequestBuilder{
        request: Request{
            Headers: make(map[string]string),
        },
    }
}

// SetMethod — метод для установки метода запроса
func (b *RequestBuilder) SetMethod(method string) *RequestBuilder {
    b.request.Method = method
    return b
}

// SetURL — метод для установки URL
func (b *RequestBuilder) SetURL(url string) *RequestBuilder {
    b.request.URL = url
    return b
}

// AddHeader — метод для добавления заголовка
func (b *RequestBuilder) AddHeader(key, value string) *RequestBuilder {
    b.request.Headers[key] = value
    return b
}

// SetTimeout — метод для установки таймаута
func (b *RequestBuilder) SetTimeout(timeout time.Duration) *RequestBuilder {
    b.request.Timeout = timeout
    return b
}

// SetBody — метод для установки тела запроса
func (b *RequestBuilder) SetBody(body string) *RequestBuilder {
    b.request.Body = body
    return b
}

// Build — конечный метод, возвращающий готовый запрос
func (b *RequestBuilder) Build() Request {
    return b.request
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "httpbuilder"
    "time"
)

func main() {
    request := httpbuilder.NewRequestBuilder().
        SetMethod("GET").
        SetURL("https://api.example.com/data").
        AddHeader("Content-Type", "application/json").
        SetTimeout(30 * time.Second).
        Build()

    fmt.Printf("Созданный запрос: %+v\n", request)
}
```

**Вывод (примерный):**
```
Созданный запрос: {Method:GET URL:https://api.example.com/data Headers:map[Content-Type:application/json] Timeout:30s Body:}
```

---

### 5.2. Построение конфигурации приложения
Создание конфигурации с опциональными полями:

```go
package configbuilder

import "time"

// Config — конечный объект конфигурации
type Config struct {
    Host     string
    Port     int
    Timeout  time.Duration
    Debug    bool
}

// ConfigBuilder — структура для пошагового построения Config
type ConfigBuilder struct {
    config Config
}

// NewConfigBuilder — конструктор для создания строителя
func NewConfigBuilder() *ConfigBuilder {
    return &ConfigBuilder{
        config: Config{
            Port: 8080, // Значение по умолчанию
        },
    }
}

// SetHost — метод для установки хоста
func (b *ConfigBuilder) SetHost(host string) *ConfigBuilder {
    b.config.Host = host
    return b
}

// SetPort — метод для установки порта
func (b *ConfigBuilder) SetPort(port int) *ConfigBuilder {
    b.config.Port = port
    return b
}

// SetTimeout — метод для установки таймаута
func (b *ConfigBuilder) SetTimeout(timeout time.Duration) *ConfigBuilder {
    b.config.Timeout = timeout
    return b
}

// SetDebug — метод для включения/выключения режима отладки
func (b *ConfigBuilder) SetDebug(debug bool) *ConfigBuilder {
    b.config.Debug = debug
    return b
}

// Build — конечный метод, возвращающий готовую конфигурацию
func (b *ConfigBuilder) Build() Config {
    return b.config
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "configbuilder"
    "time"
)

func main() {
    config := configbuilder.NewConfigBuilder().
        SetHost("localhost").
        SetPort(8081).
        SetTimeout(30 * time.Second).
        SetDebug(true).
        Build()

    fmt.Printf("Конфигурация: %+v\n", config)
}
```

**Вывод (примерный):**
```
Конфигурация: {Host:localhost Port:8081 Timeout:30s Debug:true}
```

---

## 6. Рекомендации по использованию Builder в Go

1. **Используйте для сложных объектов**: Builder оправдан, когда объект имеет множество опциональных полей или сложную логику создания.
2. **Избегайте избыточности**: Для простых объектов с небольшим числом полей используйте прямые конструкторы или структуры с опциональными параметрами.
3. **Добавляйте валидацию**: Включайте проверки данных в методы Builder для обеспечения корректности объектов.
4. **Тестирование**: Тестируйте Builder отдельно, проверяя каждый метод и конечный результат.
5. **Прозрачность**: Сделайте API Builder интуитивно понятным, используя цепочку вызовов (`*Builder` возвращает указатель на себя).

---

## 7. Преимущества и недостатки

### Преимущества:
- **Читаемость**: Позволяет создавать объекты с чёткой последовательностью вызовов.
- **Гибкость**: Поддерживает опциональные поля и сложные конфигурации.
- **Валидация**: Легко добавлять проверки данных.
- **Расширяемость**: Можно добавлять новые методы без изменения существующего кода.

### Недостатки:
- **Усложнение кода**: Для простых объектов Builder может быть избыточным.
- **Избыточность**: Требует создания дополнительной структуры (Builder), что увеличивает объём кода.
- **Сложности с отладкой**: Если Builder содержит сложную логику, отладка может стать труднее.

---

## 8. Заключение

Шаблон Builder в Go — мощный инструмент для создания сложных объектов с гибкой конфигурацией. В Go он особенно удобен благодаря поддержке структур, методов и цепочки вызовов. Используйте Builder, когда нужно поэтапно строить объекты с множеством опциональных полей или валидацией, но избегайте его для простых случаев, где достаточно прямого создания объектов.