# Работа с GORM (ORM) в Golang

## Введение

GORM — это популярная библиотека Object-Relational Mapping (ORM) для языка программирования Go. Она упрощает взаимодействие с базами данных, позволяя работать с данными как с объектами Go, вместо написания сырых SQL-запросов. GORM поддерживает множество СУБД, включая PostgreSQL, MySQL, SQLite и другие, и предоставляет мощные инструменты для CRUD-операций (Create, Read, Update, Delete), транзакций, миграций и многого другого.

В этой лекции мы разберём:
- Что такое ORM и зачем оно нужно.
- Установка и настройка GORM.
- Основные операции с GORM (CRUD).
- Работа с ассоциациями, миграциями и транзакциями.
- Практические примеры с PostgreSQL.

---

## 1. Что такое ORM?

ORM (Object-Relational Mapping) — это техника, которая позволяет сопоставлять объекты в языке программирования (например, структуры Go) с таблицами и строками в базе данных. Вместо написания сложных SQL-запросов вы работаете с объектами и методами, что ускоряет разработку и делает код более читаемым.

### Преимущества GORM:
- Упрощение работы с базой данных.
- Поддержка множества СУБД.
- Автоматическая генерация SQL-запросов.
- Обработка транзакций и миграций.
- Поддержка ассоциаций (1:1, 1:N, N:N).

### Недостатки:
- Потенциальная потеря производительности из-за автоматической генерации SQL.
- Сложность отладки для сложных запросов.
- Зависимость от библиотеки (нужно изучить её API).

---

## 2. Установка и настройка GORM

### Установка
Для начала установите GORM и драйвер для PostgreSQL:

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

### Подключение к PostgreSQL
Создадим подключение к базе данных PostgreSQL:

```go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func main() {
    dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    log.Println("Connected to PostgreSQL successfully!")
}
```

- `dsn` (Data Source Name) — строка подключения к PostgreSQL.
- `gorm.Open()` инициализирует подключение с указанным драйвером и конфигурацией.

---

## 3. Определение моделей

В GORM модели — это структуры Go, которые сопоставляются с таблицами в базе данных. Используйте теги (`gorm` и `json`) для настройки.

### Пример модели
Создадим модель `User`:

```go
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"type:varchar(100);not null"`
    Email     string `gorm:"unique;type:varchar(100)"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

- `gorm:"primaryKey"` — указывает, что `ID` — это первичный ключ.
- `gorm:"type:varchar(100);not null"` — задаёт тип и ограничение для `Name`.
- `gorm:"unique"` — делает `Email` уникальным.
- `DeletedAt` — для "мягкого" удаления (soft delete).

---

## 4. Базовые операции CRUD

### 4.1. Создание (Create)
Добавим нового пользователя:

```go
func createUser(db *gorm.DB) {
    user := User{Name: "Иван Иванов", Email: "ivan@example.com"}
    if err := db.Create(&user).Error; err != nil {
        log.Printf("Failed to create user: %v", err)
        return
    }
    log.Printf("Created user with ID: %d", user.ID)
}
```

- `db.Create()` вставляет запись в таблицу и обновляет структуру `user` с сгенерированными значениями (например, `ID`).

### 4.2. Чтение (Read)
Найдём пользователя по ID:

```go
func getUser(db *gorm.DB, id uint) {
    var user User
    if err := db.First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            log.Println("User not found")
            return
        }
        log.Printf("Failed to retrieve user: %v", err)
        return
    }
    log.Printf("Found user: %+v", user)
}
```

- `db.First()` извлекает первую запись, соответствующую условию (по `ID`).

### 4.3. Обновление (Update)
Обновим email пользователя:

```go
func updateUser(db *gorm.DB, id uint, newEmail string) {
    var user User
    if err := db.First(&user, id).Error; err != nil {
        log.Printf("User not found: %v", err)
        return
    }
    user.Email = newEmail
    if err := db.Save(&user).Error; err != nil {
        log.Printf("Failed to update user: %v", err)
        return
    }
    log.Println("User updated successfully!")
}
```

- `db.Save()` обновляет всю запись.

### 4.4. Удаление (Delete)
Удалим пользователя (или выполним "мягкое" удаление):

```go
func deleteUser(db *gorm.DB, id uint) {
    if err := db.Delete(&User{}, id).Error; err != nil {
        log.Printf("Failed to delete user: %v", err)
        return
    }
    log.Println("User deleted successfully!")
}
```

- `db.Delete()` удаляет запись. Если используется `DeletedAt`, это будет "мягкое" удаление (строка помечается, но не удаляется физически).

---

## 5. Ассоциации

GORM поддерживает отношения между таблицами: 1:1, 1:N, N:N.

### 5.1. 1:N (Один ко многим)
Добавим модель `Order` для пользователей:

```go
type Order struct {
    ID     uint   `gorm:"primaryKey"`
    Amount float64
    UserID uint
    User   User `gorm:"foreignKey:UserID"`
}
```

Связываем пользователя с заказами:

```go
func createOrderWithUser(db *gorm.DB) {
    user := User{Name: "Мария Петрова", Email: "maria@example.com"}
    if err := db.Create(&user).Error; err != nil {
        log.Printf("Failed to create user: %v", err)
        return
    }

    order := Order{Amount: 200.50, UserID: user.ID}
    if err := db.Create(&order).Error; err != nil {
        log.Printf("Failed to create order: %v", err)
        return
    }
    log.Println("Order created successfully!")
}
```

Получим пользователя с его заказами:

```go
func getUserWithOrders(db *gorm.DB, id uint) {
    var user User
    if err := db.Preload("Order").First(&user, id).Error; err != nil {
        log.Printf("Failed to retrieve user: %v", err)
        return
    }
    log.Printf("User: %+v, Orders: %+v", user, user.Order)
}
```

- `Preload("Order")` загружает связанные записи (eager loading).

### 5.2. N:N (Многие ко многим)
Добавим связь между пользователями и ролями:

```go
type Role struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"type:varchar(50)"`
    Users []User `gorm:"many2many:user_roles;"`
}

type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"type:varchar(100);not null"`
    Email     string `gorm:"unique;type:varchar(100)"`
    Roles     []Role `gorm:"many2many:user_roles;"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

Создаём промежуточную таблицу `user_roles` автоматически:

```go
func createUserWithRoles(db *gorm.DB) {
    user := User{Name: "Алексей Сидоров", Email: "alexey@example.com"}
    role := Role{Name: "Admin"}
    user.Roles = []Role{role}

    if err := db.Create(&user).Error; err != nil {
        log.Printf("Failed to create user and roles: %v", err)
        return
    }
    log.Println("User with roles created successfully!")
}
```

---

## 6. Миграции

GORM позволяет автоматически создавать или обновлять таблицы на основе моделей:

```go
func migrate(db *gorm.DB) {
    if err := db.AutoMigrate(&User{}, &Order{}, &Role{}); err != nil {
        log.Printf("Failed to migrate database: %v", err)
        return
    }
    log.Println("Database migrated successfully!")
}
```

- `AutoMigrate()` создаёт таблицы, если их нет, или изменяет существующие, подстраивая под модель.

---

## 7. Транзакции

GORM поддерживает транзакции для атомарных операций:

```go
func transferMoney(db *gorm.DB, fromID, toID uint, amount float64) error {
    tx := db.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    var fromUser, toUser User
    if err := tx.First(&fromUser, fromID).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.First(&toUser, toID).Error; err != nil {
        tx.Rollback()
        return err
    }

    fromUser.Balance -= amount
    toUser.Balance += amount

    if err := tx.Save(&fromUser).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Save(&toUser).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
```

- `db.Begin()` начинает транзакцию.
- `tx.Commit()` фиксирует изменения.
- `tx.Rollback()` откатывает при ошибке.

---

## 8. Дополнительные возможности

### 8.1. Условия и фильтры
GORM позволяет строить сложные запросы:

```go
func findUsersByName(db *gorm.DB, name string) {
    var users []User
    if err := db.Where("name LIKE ?", "%"+name+"%").Find(&users).Error; err != nil {
        log.Printf("Failed to find users: %v", err)
        return
    }
    log.Printf("Found users: %+v", users)
}
```

- `Where()` добавляет условие фильтрации.

### 8.2. Пагинация и сортировка
```go
func getPaginatedUsers(db *gorm.DB, page, pageSize int) {
    var users []User
    if err := db.Order("created_at desc").Limit(pageSize).Offset((page-1)*pageSize).Find(&users).Error; err != nil {
        log.Printf("Failed to retrieve users: %v", err)
        return
    }
    log.Printf("Paginated users: %+v", users)
}
```

- `Order()`, `Limit()`, `Offset()` — для сортировки и пагинации.

---

## 9. Полезные советы

1. **Производительность**: Для сложных запросов используйте сырые SQL-запросы через `db.Raw()` или `db.Exec()`, если GORM не оптимален.
2. **Логирование**: Включите логирование для отладки: `db = db.Debug()`.
3. **Тестирование**: Используйте `sqlite` для тестов, так как оно быстрее и проще настраивается.
4. **Миграции**: Используйте `AutoMigrate()` осторожно в продакшене, так как это может привести к потере данных.

---

## 10. Заключение

GORM — мощный инструмент для работы с базами данных в Go. Он упрощает CRUD-операции, миграции и управление ассоциациями, но требует понимания его возможностей и ограничений. Освоив GORM, вы сможете быстрее разрабатывать приложения, сохраняя читаемость и безопасность кода.