# Шаблон проектирования Factory Method в Golang

## Введение

Шаблон проектирования **Factory Method** (Фабричный метод) — это порождающий шаблон, который определяет интерфейс для создания объектов, но позволяет подклассам (или реализациям) решать, какой конкретный класс создавать. В Go, где нет классов в традиционном смысле, этот шаблон реализуется с использованием интерфейсов, структур и функций, что делает его особенно подходящим для гибкого создания объектов.

В этой лекции мы разберём:
- Что такое Factory Method и где он применяется.
- Как реализовать Factory Method в Go.
- Преимущества и недостатки шаблона.
- Примеры использования в реальных задачах.
- Рекомендации по применению в Go.

---

## 1. Что такое Factory Method?

Factory Method — это шаблон, который:
- Предоставляет интерфейс для создания объектов, но делегирует создание конкретных объектов подклассам или отдельным функциям.
- Избегает жёсткой привязки к конкретным классам, улучшая гибкость и расширяемость кода.
- Используется, когда нужно создавать объекты без указания их конкретного типа на этапе компиляции.

### Примеры использования:
- Создание различных типов документов (PDF, Word, HTML) в зависимости от входных данных.
- Фабрика для создания объектов базы данных (PostgreSQL, MySQL, SQLite).
- Генерация различных форматов логов (файл, консоль, сеть).

---

## 2. Реализация Factory Method в Go

В Go Factory Method реализуется через интерфейсы и функции, которые возвращают объекты, соответствующие этому интерфейсу. Go не использует классы или наследование, поэтому мы опираемся на композицию и интерфейсы.

### 2.1. Базовая структура

Рассмотрим пример фабрики для создания транспортных средств (Vehicle). Мы определим интерфейс `Vehicle`, конкретные типы (например, `Car` и `Bike`) и фабрику для их создания.

#### Шаг 1: Определение интерфейса
```go
package factory

// Vehicle — интерфейс для транспортных средств
type Vehicle interface {
    Drive() string
}
```

#### Шаг 2: Конкретные реализации
```go
// Car — конкретная реализация транспортного средства (автомобиль)
type Car struct{}

func (c *Car) Drive() string {
    return "Машина едет по дороге!"
}

// Bike — конкретная реализация транспортного средства (велосипед)
type Bike struct{}

func (b *Bike) Drive() string {
    return "Велосипед едет по тропинке!"
}
```

#### Шаг 3: Фабрика (Factory Method)
Создадим функцию, которая будет создавать объекты `Vehicle` в зависимости от типа:

```go
// CreateVehicle — фабричный метод для создания транспортных средств
func CreateVehicle(vehicleType string) Vehicle {
    switch vehicleType {
    case "car":
        return &Car{}
    case "bike":
        return &Bike{}
    default:
        return nil
    }
}
```

#### Шаг 4: Использование
```go
package main

import (
    "fmt"
    "factory"
)

func main() {
    // Создаём автомобиль
    car := factory.CreateVehicle("car")
    if car != nil {
        fmt.Println(car.Drive()) // Вывод: Машина едет по дороге!
    }

    // Создаём велосипед
    bike := factory.CreateVehicle("bike")
    if bike != nil {
        fmt.Println(bike.Drive()) // Вывод: Велосипед едет по тропинке!
    }

    // Неправильный тип
    unknown := factory.CreateVehicle("truck")
    if unknown == nil {
        fmt.Println("Неизвестный тип транспортного средства")
    }
}
```

**Вывод:**
```
Машина едет по дороге!
Велосипед едет по тропинке!
Неизвестный тип транспортного средства
```

---

### 2.2. Расширенная реализация с конфигурацией

Добавим возможность передавать конфигурацию для создания объектов, чтобы сделать фабрику более гибкой:

```go
package factory

// VehicleConfig — конфигурация для создания транспортных средств
type VehicleConfig struct {
    Type     string
    MaxSpeed int
}

// Vehicle — интерфейс для транспортных средств
type Vehicle interface {
    Drive() string
    GetMaxSpeed() int
}

// Car — конкретная реализация транспортного средства (автомобиль)
type Car struct {
    maxSpeed int
}

func (c *Car) Drive() string {
    return fmt.Sprintf("Машина едет по дороге со скоростью до %d км/ч!", c.maxSpeed)
}

func (c *Car) GetMaxSpeed() int {
    return c.maxSpeed
}

// Bike — конкретная реализация транспортного средства (велосипед)
type Bike struct {
    maxSpeed int
}

func (b *Bike) Drive() string {
    return fmt.Sprintf("Велосипед едет по тропинке со скоростью до %d км/ч!", b.maxSpeed)
}

func (b *Bike) GetMaxSpeed() int {
    return b.maxSpeed
}

// CreateVehicle — фабричный метод для создания транспортных средств с конфигурацией
func CreateVehicle(config VehicleConfig) Vehicle {
    switch config.Type {
    case "car":
        return &Car{maxSpeed: config.MaxSpeed}
    case "bike":
        return &Bike{maxSpeed: config.MaxSpeed}
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
    "factory"
)

func main() {
    carConfig := factory.VehicleConfig{Type: "car", MaxSpeed: 120}
    car := factory.CreateVehicle(carConfig)
    if car != nil {
        fmt.Println(car.Drive()) // Вывод: Машина едет по дороге со скоростью до 120 км/ч!
        fmt.Printf("Максимальная скорость: %d км/ч\n", car.GetMaxSpeed())
    }

    bikeConfig := factory.VehicleConfig{Type: "bike", MaxSpeed: 30}
    bike := factory.CreateVehicle(bikeConfig)
    if bike != nil {
        fmt.Println(bike.Drive()) // Вывод: Велосипед едет по тропинке со скоростью до 30 км/ч!
        fmt.Printf("Максимальная скорость: %d км/ч\n", bike.GetMaxSpeed())
    }
}
```

**Вывод:**
```
Машина едет по дороге со скоростью до 120 км/ч!
Максимальная скорость: 120 км/ч
Велосипед едет по тропинке со скоростью до 30 км/ч!
Максимальная скорость: 30 км/ч
```

---

## 3. Преимущества Factory Method

- **Гибкость**: Позволяет создавать объекты разных типов без изменения клиентского кода, добавляя новые типы через новые реализации интерфейса.
- **Инкапсуляция**: Скрывает логику создания объектов, упрощая клиентский код.
- **Расширяемость**: Легко добавлять новые типы объектов, не нарушая существующий код.
- **Соответствие принципам OOP**: Следует принципу открытости/закрытости (Open/Closed Principle).

---

## 4. Недостатки Factory Method

- **Усложнение кода**: Добавление фабрики может увеличить сложность, особенно если создаваемых типов много.
- **Избыточность**: Для простых случаев может быть избыточным, когда достаточно прямого создания объектов.
- **Тестирование**: Может быть трудно замокать или тестировать фабрику, если она сильно интегрирована с логикой приложения.

---

## 5. Примеры реального использования

### 5.1. Фабрика подключений к базам данных
Создание подключений к разным СУБД (PostgreSQL, MySQL):

```go
package database

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// Database — интерфейс для работы с базой данных
type Database interface {
    Connect() error
    Query(string) error
}

// PostgresDB — реализация для PostgreSQL
type PostgresDB struct {
    db *gorm.DB
}

func (p *PostgresDB) Connect() error {
    dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
    }
    p.db = db
    return nil
}

func (p *PostgresDB) Query(query string) error {
    return p.db.Raw(query).Error
}

// MySQLDB — реализация для MySQL
type MySQLDB struct {
    db *gorm.DB
}

func (m *MySQLDB) Connect() error {
    dsn := "user:password@/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to MySQL: %v", err)
    }
    m.db = db
    return nil
}

func (m *MySQLDB) Query(query string) error {
    return m.db.Raw(query).Error
}

// CreateDatabase — фабричный метод для создания подключений
func CreateDatabase(dbType string) Database {
    switch dbType {
    case "postgres":
        return &PostgresDB{}
    case "mysql":
        return &MySQLDB{}
    default:
        return nil
    }
}
```

#### Использование:
```go
package main

import (
    "database"
    "fmt"
)

func main() {
    // Подключение к PostgreSQL
    pgDB := database.CreateDatabase("postgres")
    if pgDB != nil {
        if err := pgDB.Connect(); err != nil {
            fmt.Printf("Ошибка подключения: %v\n", err)
        } else {
            fmt.Println("Подключение к PostgreSQL успешно!")
            pgDB.Query("SELECT * FROM users")
        }
    }

    // Подключение к MySQL
    mysqlDB := database.CreateDatabase("mysql")
    if mysqlDB != nil {
        if err := mysqlDB.Connect(); err != nil {
            fmt.Printf("Ошибка подключения: %v\n", err)
        } else {
            fmt.Println("Подключение к MySQL успешно!")
            mysqlDB.Query("SELECT * FROM users")
        }
    }
}
```

---

### 5.2. Фабрика форматов логов
Создание различных логгеров (консоль, файл):

```go
package logger

import "fmt"

// Logger — интерфейс для логгеров
type Logger interface {
    Log(message string)
}

// ConsoleLogger — логгер для вывода в консоль
type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(message string) {
    fmt.Println("CONSOLE:", message)
}

// FileLogger — логгер для записи в файл (простая имитация)
type FileLogger struct{}

func (f *FileLogger) Log(message string) {
    fmt.Println("FILE:", message)
}

// CreateLogger — фабричный метод для создания логгеров
func CreateLogger(loggerType string) Logger {
    switch loggerType {
    case "console":
        return &ConsoleLogger{}
    case "file":
        return &FileLogger{}
    default:
        return nil
    }
}
```

#### Использование:
```go
package main

import (
    "logger"
)

func main() {
    consoleLogger := logger.CreateLogger("console")
    if consoleLogger != nil {
        consoleLogger.Log("Приложение запущено")
    }

    fileLogger := logger.CreateLogger("file")
    if fileLogger != nil {
        fileLogger.Log("Ошибка в приложении")
    }
}
```

**Вывод:**
```
CONSOLE: Приложение запущено
FILE: Ошибка в приложении
```

---

## 6. Рекомендации по использованию Factory Method в Go

1. **Используйте интерфейсы**: Определите интерфейс для создаваемых объектов, чтобы обеспечить гибкость и расширяемость.
2. **Избегайте избыточности**: Для простых случаев прямое создание объектов может быть проще, чем создание фабрики.
3. **Тестирование**: Создавайте мок-объекты для фабрик, чтобы легко тестировать клиентский код.
4. **Конфигурация**: Передавайте конфигурацию через структуры, чтобы сделать фабрику более гибкой.
5. **Минимизируйте зависимости**: Убедитесь, что фабрика не зависит от слишком многих внешних компонентов.

---

## 7. Преимущества и недостатки

### Преимущества:
- **Гибкость**: Легко добавлять новые типы объектов без изменения клиентского кода.
- **Инкапсуляция**: Скрывает сложность создания объектов.
- **Соответствие принципам SOLID**: Следует принципу открытости/закрытости.

### Недостатки:
- **Усложнение кода**: Добавление фабрики может сделать код менее читаемым для простых случаев.
- **Избыточность**: Может быть избыточным для небольших приложений.
- **Сложности с отладкой**: Трудно отследить, какие объекты создаются, если фабрика сложная.

---

## 8. Заключение

Шаблон Factory Method в Go — мощный инструмент для создания объектов с гибкостью и расширяемостью. В Go он особенно удобен благодаря поддержке интерфейсов и функций. Используйте этот шаблон, когда нужно изолировать логику создания объектов и поддерживать открытость для новых типов. Однако для простых случаев предпочтительнее избегать избыточного использования фабрик, полагаясь на прямое создание объектов или передачу зависимостей.