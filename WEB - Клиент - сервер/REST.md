# REST

## Введение

**REST** (Representational State Transfer) — это архитектурный стиль для создания масштабируемых веб-сервисов. Он был предложен Роем Филдингом в его диссертации 2000 года и стал стандартом де-факто для API в веб-разработке. REST не является ни протоколом, ни библиотекой, а набором принципов проектирования, которые делают системы простыми, предсказуемыми и совместимыми.

Эта лекция объяснит ключевые концепции REST, его принципы (ограничения), как реализовать REST API на Go, а также разберёт типичные ошибки и лучшие практики для разработчиков.

---

## Что такое REST?

REST — это способ организации взаимодействия между клиентом и сервером через HTTP. Основная идея: ресурсы (например, пользователи, заказы) представляются в виде URI, а операции с ними выполняются через стандартные методы HTTP (GET, POST, PUT, DELETE).

### Основные характеристики REST

- **Ресурсы**: Всё в REST — это ресурс (сущность), идентифицируемый уникальным URI (например, `/users/123`).
- **Клиент-сервер**: Разделение логики клиента (UI) и сервера (данные).
- **Без состояния (Stateless)**: Каждый запрос содержит всю информацию, необходимую для его обработки.
- **Кэшируемость**: Ответы сервера могут кэшироваться для повышения производительности.
- **Единообразие интерфейса**: Стандартизированный способ взаимодействия с ресурсами.

---

## Принципы REST (ограничения Филдинга)

REST базируется на шести принципах, которые Senior-разработчик должен понимать и применять:

1. **Клиент-сервер**: Разделение ответственности между клиентом и сервером улучшает независимость и масштабируемость.
2. **Stateless (Без состояния)**: Сервер не хранит информацию о состоянии клиента между запросами. Каждый запрос самодостаточен.
3. **Cacheable (Кэшируемость)**: Клиент может кэшировать ответы, если сервер указывает это (например, через заголовки `Cache-Control`).
4. **Uniform Interface (Единообразный интерфейс)**: 
   - Ресурсы идентифицируются через URI.
   - Стандартные методы HTTP (GET, POST, PUT, DELETE).
   - HATEOAS (Hypermedia as the Engine of Application State) — опционально, ссылки в ответах для навигации.
5. **Layered System (Слоистая система)**: Клиент не знает, общается он напрямую с сервером или через прокси/балансировщик.
6. **Code on Demand (Код по запросу)**: Опционально, сервер может отправлять исполняемый код (например, JavaScript). В API редко используется.

---

## HTTP-методы в REST

REST использует стандартные методы HTTP для работы с ресурсами:

| Метод   | Действие                | Пример URI         | Описание                            |
|---------|-------------------------|--------------------|-------------------------------------|
| GET     | Получение ресурса       | `/users/123`       | Возвращает данные пользователя 123  |
| POST    | Создание ресурса        | `/users`           | Создаёт нового пользователя         |
| PUT     | Обновление ресурса      | `/users/123`       | Обновляет данные пользователя 123   |
| DELETE  | Удаление ресурса        | `/users/123`       | Удаляет пользователя 123            |
| PATCH   | Частичное обновление    | `/users/123`       | Изменяет часть данных пользователя  |

---

## Реализация REST API на Go

Давайте создадим простой REST API на Go для управления пользователями. Мы будем использовать стандартную библиотеку `net/http`, чтобы показать базовые принципы, а затем добавим улучшения.

### Пример: Базовый REST API

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// User представляет модель пользователя
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Хранилище пользователей (в памяти)
var users = make(map[string]User)

func main() {
	http.HandleFunc("/users", usersHandler)       // Список и создание
	http.HandleFunc("/users/", userHandler)       // Получение, обновление, удаление
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

// usersHandler обрабатывает запросы к /users
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Получение списка пользователей
		json.NewEncoder(w).Encode(users)
	case http.MethodPost:
		// Создание нового пользователя
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		users[user.ID] = user
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// userHandler обрабатывает запросы к /users/{id}
func userHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Получение пользователя
		user, exists := users[id]
		if !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	case http.MethodPut:
		// Обновление пользователя
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		user.ID = id
		users[id] = user
		json.NewEncoder(w).Encode(user)
	case http.MethodDelete:
		// Удаление пользователя
		if _, exists := users[id]; !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		delete(users, id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
```

#### Как запустить?
- Сохраните код в main.go.
- Выполните go run main.go.
- Используйте curl или Postman для тестирования:
- curl -X GET http://localhost:8080/users — получить всех пользователей.
- curl -X POST -d '{"id":"1","name":"Alice","email":"alice@example.com"}' http://localhost:8080/users — создать пользователя.
- curl -X GET http://localhost:8080/users/1 — получить пользователя с ID 1.
- curl -X DELETE http://localhost:8080/users/1 — удалить пользователя.

#### Объяснение
- Маршруты: /users для списка и создания, /users/{id} для операций с конкретным пользователем.
- Методы: Обрабатываем GET, POST, PUT, DELETE с соответствующими статусами HTTP (200, 201, 204, 404, 405).
- Хранилище: Используем map в памяти для простоты.

## Лучшие практики для REST API
### Используйте правильные статусы HTTP:
  - 200 OK — успешный GET или PUT.
  - 201 Created — успешный POST.
  - 204 No Content — успешный DELETE.
  - 400 Bad Request — ошибка в запросе.
  - 404 Not Found — ресурс не найден.
  - 500 Internal Server Error — ошибка сервера.
### Версионирование API:
  - Используйте /v1/users вместо /users, чтобы поддерживать обратную совместимость.
  - HATEOAS (опционально):

### Добавляйте ссылки в ответы:
```
{
  "id": "1",
  "name": "Alice",
  "links": {
    "self": "/users/1",
    "delete": "/users/1"
  }
}
```
### Кэширование:
Используйте заголовки ETag или Cache-Control для оптимизации.

### Валидация данных:
  - Проверяйте входные данные перед обработкой.

### Безопасность:
  - Используйте HTTPS.
  - Добавьте аутентификацию (JWT, OAuth).

 Типичные ошибки
  - Нарушение Stateless: Хранение сессий на сервере вместо передачи состояния в запросе.
  - Неправильные методы: Использование GET для изменения данных (должно быть POST/PUT).
  - Сложные URI: /users/getById/123 вместо /users/123.
  - Игнорирование статусов: Возврат 200 для всех ответов вместо 404 или 400.
  - Отсутствие документации: REST API должен быть легко читаемым и документированным (например, через OpenAPI/Swagger).

## Заключение
REST — это не просто "HTTP + JSON", а архитектурный стиль, требующий соблюдения строгих принципов. Pазработчик должен уметь проектировать API, которые:
  - Масштабируемы и поддерживаемы.
  - Соответствуют стандартам REST.
  - Удобны для клиентов.
  - Go отлично подходит для REST благодаря своей простоте, производительности и стандартной библиотеке net/http. Однако для сложных проектов стоит использовать фреймворки (Gin, Echo) или библиотеки вроде gorilla/mux.

Совет: Начинайте с минимального API и добавляйте сложности (версионирование, кэширование) по мере необходимости. REST — это про баланс между простотой и функциональностью.
