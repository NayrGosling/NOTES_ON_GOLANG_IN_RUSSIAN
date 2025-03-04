# Шаблон проектирования Singleton в Golang

## Введение

Шаблон проектирования **Singleton** (Одиночка) — это порождающий шаблон, который гарантирует, что класс имеет только один экземпляр и предоставляет глобальную точку доступа к этому экземпляру. В Go, где нет классов в традиционном смысле, Singleton реализуется с использованием структур, функций и принципов языка, таких как инициализация пакетов и атомарные операции.

В этой лекции мы разберём:
- Что такое Singleton и где он применяется.
- Как реализовать Singleton в Go.
- Преимущества и недостатки шаблона.
- Примеры использования в реальных задачах.
- Рекомендации по применению в Go.

---

## 1. Что такое Singleton?

Singleton — это шаблон, который:
- Обеспечивает существование только одного экземпляра класса (или структуры в Go).
- Предоставляет глобальный доступ к этому экземпляру через статический метод или функцию.

### Примеры использования:
- Логирование: один логгер для всего приложения.
- Настройки: единый объект для хранения глобальных конфигураций.
- Подключения к базе данных: одно подключение для оптимизации ресурсов.

---

## 2. Реализация Singleton в Go

В Go нет встроенных классов, но мы можем использовать структуры, функции и пакеты для реализации Singleton. Основные подходы включают:
- Использование инициализации пакета `init`.
- Ленивую инициализацию с синхронизацией для потокобезопасности.
- Использование пакета `sync` для атомарных операций.

### 2.1. Простая реализация без потокобезопасности

Рассмотрим простой пример Singleton для объекта конфигурации:

```go
package singleton

import "sync"

// Config — структура для хранения конфигурации
type Config struct {
    Host string
    Port int
}

// instance — приватная переменная для хранения единственного экземпляра
var instance *Config

// GetInstance возвращает единственный экземпляр Config
func GetInstance() *Config {
    if instance == nil {
        instance = &Config{
            Host: "localhost",
            Port: 8080,
        }
    }
    return instance
}
```

#### Использование:
```go
package main

import (
    "fmt"
    "singleton"
)

func main() {
    config1 := singleton.GetInstance()
    config2 := singleton.GetInstance()

    fmt.Printf("Config1: %v\n", config1)
    fmt.Printf("Config2: %v\n", config2)
    fmt.Printf("Один и тот же объект: %t\n", config1 == config2)
}
```

**Вывод:**
```
Config1: &{localhost 8080}
Config2: &{localhost 8080}
Один и тот же объект: true
```

#### Ограничения:
- Эта реализация не потокобезопасна. Если несколько горутин одновременно вызовут `GetInstance()`, могут создаться несколько экземпляров.

---

### 2.2. Потокобезопасная реализация с `sync.Once`

В Go для потокобезопасной ленивой инициализации часто используется `sync.Once`, который гарантирует, что инициализация выполнится только один раз, даже при одновременных вызовах из разных горутин.

```go
package singleton

import "sync"

// Config — структура для хранения конфигурации
type Config struct {
    Host string
    Port int
}

// instance — приватная переменная для хранения единственного экземпляра
var instance *Config
var once sync.Once

// GetInstance возвращает единственный экземпляр Config
func GetInstance() *Config {
    once.Do(func() {
        instance = &Config{
            Host: "localhost",
            Port: 8080,
        }
    })
    return instance
}
```

#### Преимущества:
- `sync.Once` использует мьютекс для синхронизации, предотвращая создание дубликатов в многопоточной среде.
- Выполнение инициализации происходит только один раз.

#### Использование:
Тот же код в `main`, но теперь он безопасен для многопоточных приложений:

```go
package main

import (
    "fmt"
    "singleton"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            config := singleton.GetInstance()
            fmt.Printf("Горутина %d: Config: %v\n", i, config)
        }()
    }
    wg.Wait()
    // Все горутины получат один и тот же экземпляр
}
```

**Вывод (примерный, порядок может варьироваться):**
```
Горутина 0: Config: &{localhost 8080}
Горутина 2: Config: &{localhost 8080}
Горутина 1: Config: &{localhost 8080}
Горутина 4: Config: &{localhost 8080}
Горутина 3: Config: &{localhost 8080}
```

---

### 2.3. Инициализация через `init`

В Go можно использовать функцию `init` для создания Singleton на этапе инициализации пакета:

```go
package singleton

import "sync"

// Config — структура для хранения конфигурации
type Config struct {
    Host string
    Port int
}

var instance *Config
var mutex sync.Mutex

func init() {
    instance = &Config{
        Host: "localhost",
        Port: 8080,
    }
}

// GetInstance возвращает единственный экземпляр Config
func GetInstance() *Config {
    mutex.Lock()
    defer mutex.Unlock()
    return instance
}
```

#### Преимущества:
- Инициализация происходит один раз при загрузке пакета.
- Подходит для случаев, когда Singleton не требует ленивой загрузки.

#### Ограничения:
- Не ленивая инициализация — экземпляр создаётся даже если он не используется.
- Требует ручной синхронизации (`mutex`), если нужно обновлять данные.

---

## 3. Преимущества Singleton

- **Глобальный доступ**: Предоставляет единую точку доступа к ресурсу.
- **Экономия ресурсов**: Один экземпляр вместо множества, что особенно важно для подключений (например, к базе данных или логгерам).
- **Упрощение кода**: Избегает дублирования логики для создания экземпляров.

---

## 4. Недостатки Singleton

- **Нарушение принципа единственной ответственности**: Singleton может стать "мусоросборником" для глобального состояния.
- **Сложности с тестированием**: Глобальный объект трудно подменить мок-объектом.
- **Потенциальные проблемы с потокобезопасностью**: Если не синхронизировать доступ, могут возникнуть гонки данных.
- **Скрытая зависимость**: Код, зависящий от Singleton, может стать менее гибким и трудным для поддержки.

---

## 5. Примеры реального использования

### 5.1. Логгер
Singleton идеально подходит для логгера, так как приложение обычно использует один логгер:

```go
package logger

import (
    "log"
    "sync"
)

type Logger struct {
    logger *log.Logger
}

var instance *Logger
var once sync.Once

func GetInstance() *Logger {
    once.Do(func() {
        instance = &Logger{
            logger: log.Default(),
        }
    })
    return instance
}

func (l *Logger) Info(msg string) {
    l.logger.Println("INFO:", msg)
}
```

#### Использование:
```go
package main

import (
    "logger"
)

func main() {
    logger := logger.GetInstance()
    logger.Info("Приложение запущено")
}
```

### 5.2. Подключение к базе данных
Singleton для пула подключений:

```go
package db

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "sync"
)

var instance *gorm.DB
var once sync.Once

func GetInstance() *gorm.DB {
    once.Do(func() {
        dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            panic("Failed to connect to database: " + err.Error())
        }
        instance = db
    })
    return instance
}
```

#### Использование:
```go
package main

import (
    "db"
    "log"
)

func main() {
    db := db.GetInstance()
    if err := db.AutoMigrate(&User{}); err != nil {
        log.Fatal("Migration failed:", err)
    }
}
```

---

## 6. Альтернативы Singleton в Go

В Go часто избегают Singleton из-за его потенциальных проблем. Альтернативы включают:
- **Передача зависимостей**: Передавайте экземпляры через параметры функций или структуры.
- **Контекст (Context)**: Используйте `context.Context` для передачи глобальных данных.
- **Пакетные переменные**: Используйте глобальные переменные пакета с осторожностью и синхронизацией.

Пример передачи зависимости:

```go
package main

type Config struct {
    Host string
    Port int
}

type Service struct {
    config *Config
}

func NewService(config *Config) *Service {
    return &Service{config: config}
}

func main() {
    config := &Config{Host: "localhost", Port: 8080}
    service := NewService(config)
    // Используем service...
}
```

---

## 7. Рекомендации по использованию Singleton в Go

1. **Используйте с осторожностью**: Singleton может усложнить тестирование и поддержку. Рассмотрите передачу зависимостей как альтернативу.
2. **Обеспечьте потокобезопасность**: Всегда используйте `sync.Once` или мьютексы для многопоточных приложений.
3. **Избегайте избыточного состояния**: Не превращайте Singleton в контейнер для всех глобальных данных.
4. **Тестирование**: Для тестов создавайте мок-объекты или используйте интерфейсы для подмены Singleton.
5. **Логирование и подключения**: Используйте Singleton для логгеров, пулов подключений или других ресурсов, где единый экземпляр оправдан.

---

## 8. Заключение

Шаблон Singleton в Go — полезный инструмент для создания единственного экземпляра объекта с глобальным доступом, но он требует осторожного использования. В Go предпочтительнее передавать зависимости явно, чтобы минимизировать побочные эффекты и улучшить тестируемость. Если вы выбираете Singleton, используйте `sync.Once` для потокобезопасности и тщательно тестируйте код.