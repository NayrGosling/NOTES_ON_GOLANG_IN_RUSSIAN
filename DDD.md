# Domain-Driven Design (DDD) в Golang

## Введение

**Domain-Driven Design (DDD)** — это подход к проектированию программного обеспечения, при котором основное внимание уделяется **бизнес-логике** и предметной области. DDD помогает разрабатывать гибкие и масштабируемые приложения, структурируя код на основе реальных бизнес-процессов.

В этом материале мы рассмотрим:

- Основные концепции DDD
- Архитектурные принципы
- Реализацию DDD в Go

---

## Основные концепции DDD

### 1️⃣ **D****omain (Дом****ен)**

Это основная часть системы, включающая бизнес-логику.

### 2️⃣ **Entity (Сущность)**

Объект, имеющий **уникальный идентификатор** и определяющий бизнес-логику.

```go
type User struct {
    ID   int
    Name string
    Age  int
}
```

### 3️⃣ **Value Object (Объект-значение)**

Объект **без идентификатора**, полностью определяемый своими свойствами.

```go
type Address struct {
    City   string
    Street string
}
```

### 4️⃣ **Aggregate (Агрегат)**

Группа сущностей, управляемая одной главной сущностью.

```go
type Order struct {
    ID      int
    Items   []Item
    Total   float64
}
```

### 5️⃣ **Repository (Репозиторий)**

Отвечает за получение и сохранение сущностей.

```go
type UserRepository interface {
    Save(user User) error
    FindByID(id int) (User, error)
}
```

### 6️⃣ **Service (Сервис)**

Описывает операции, которые нельзя отнести к сущности.

```go
type UserService struct {
    repo UserRepository
}

func (s *UserService) RegisterUser(name string, age int) error {
    user := User{Name: name, Age: age}
    return s.repo.Save(user)
}
```

---

## Архитектурные принципы

### 1️⃣ **Разделение слоев**

DDD предполагает **четкое разделение слоев**:

- **Домен (Бизнес-логика)** → сущности, агрегаты, объекты-значения
- **Приложение** → сервисы, кейсы использования
- **Инфраструктура** → репозитории, API, БД

### 2️⃣ **Dependency Inversion (Инверсия зависимостей)**

Бизнес-логика не зависит от внешних деталей (например, базы данных).

```go
type OrderRepository interface {
    Save(order Order) error
}
```

Инфраструктурная реализация:

```go
type OrderRepositoryDB struct {
    db *sql.DB
}

func (r *OrderRepositoryDB) Save(order Order) error {
    _, err := r.db.Exec("INSERT INTO orders VALUES (?, ?)", order.ID, order.Total)
    return err
}
```

---

## Реализация DDD в Go

### 📌 **Структура проекта**

```
project/
├── domain/
│   ├── user.go
│   ├── order.go
│   ├── repository.go
│
├── application/
│   ├── user_service.go
│   ├── order_service.go
│
├── infrastructure/
│   ├── persistence/
│   │   ├── user_repo_db.go
│   │   ├── order_repo_db.go
│
├── main.go
```

### 📌 **Пример использования**

```go
func main() {
    db, _ := sql.Open("mysql", "user:password@/dbname")
    userRepo := &UserRepositoryDB{db: db}
    userService := UserService{repo: userRepo}
    
    err := userService.RegisterUser("Alice", 25)
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## Итог

✅ DDD помогает структурировать код и выделять бизнес-логику.
✅ В Go можно реализовать DDD через интерфейсы и инверсию зависимостей.
✅ Четкая архитектура DDD делает код **масштабируемым и поддерживаемым**.

Использование DDD в Go улучшает структуру приложения и снижает технический долг. 🚀

