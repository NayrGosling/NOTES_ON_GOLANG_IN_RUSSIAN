# WebSockets

## Введение
Привет! Сегодня мы поговорим про **WebSockets** – супербыстрый способ общения между браузером и сервером. Это как волшебный телефон, который всегда остаётся на связи, вместо того чтобы отправлять письма (HTTP-запросы).

---

## Что такое WebSockets и зачем они нужны?
Обычные HTTP-запросы – это как письма: браузер отправляет запрос 📜, сервер отвечает ✉️, и на этом всё. Чтобы получить новые данные, браузер должен снова спросить. А WebSockets – это как телефонная линия 📞: соединение открыто, и сервер может сам присылать обновления, когда нужно!

### Где применяются WebSockets?
- 🔥 **Чаты и мессенджеры** (WhatsApp, Telegram Web)
- 🎮 **Онлайн-игры** (Counter-Strike, Fortnite)
- 📈 **Торговые платформы** (биржи, котировки акций)
- 🔄 **Обновления в реальном времени** (уведомления, live-комментарии, потоковые данные)
- 🚗 **Интернет вещей (IoT)** (умные устройства, датчики)

---

## Как работают WebSockets?
WebSocket-соединение проходит три этапа:
1. 📞 **Установление соединения** (Handshake) – клиент отправляет HTTP-запрос с `Upgrade: websocket`, сервер подтверждает.
2. 🔄 **Двусторонний обмен данными** – сервер и клиент могут отправлять сообщения в любое время.
3. 🔚 **Закрытие соединения** – если соединение больше не нужно, его можно закрыть.

### HTTP vs WebSockets
| Функция            | HTTP              | WebSockets        |
|-------------------|-----------------|-----------------|
| Соединение       | Новое на каждый запрос | Одно постоянное |
| Направление данных | Только клиент → сервер | Двусторонний |
| Скорость        | Медленный (много запросов) | Быстрый (меньше накладных расходов) |
| Использование   | API-запросы, загрузка страниц | Чаты, игры, обновления в реальном времени |

---

## WebSocket Handshake (Установление соединения)
Все начинается с обычного HTTP-запроса, но с небольшим трюком:
```http
GET /chat HTTP/1.1
Host: example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13
```

Сервер отвечает:
```http
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```
Теперь соединение установлено, и можно отправлять сообщения без лишних заголовков! 🚀

---

## Как реализовать WebSockets в Go
В Go есть замечательная библиотека `github.com/gorilla/websocket`. Давай создадим простой WebSocket-сервер!

### 1. Устанавливаем библиотеку
```sh
go get -u github.com/gorilla/websocket
```

### 2. Пишем сервер
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при обновлении соединения:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Ошибка при чтении сообщения:", err)
			break
		}
		fmt.Println("Получено сообщение:", string(msg))
		conn.WriteMessage(messageType, msg) // Отправляем обратно (echo)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
```

Теперь, если открыть `ws://localhost:8080/ws` в WebSocket-клиенте, можно отправлять и получать сообщения!

---

## WebSocket-клиент (JavaScript)
```js
const socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = () => {
    console.log("Соединение установлено!");
    socket.send("Привет, сервер!");
};

socket.onmessage = (event) => {
    console.log("Получено сообщение:", event.data);
};

socket.onclose = () => {
    console.log("Соединение закрыто");
};
```
Теперь браузер и сервер могут обмениваться сообщениями в реальном времени! 🚀

---

## Масштабирование WebSockets
В реальных проектах WebSockets могут обслуживать тысячи соединений. Как с этим справляться?

- **Goroutines** – Go отлично подходит для многопоточных WebSockets!
- **Redis Pub/Sub** – для синхронизации между серверами.
- **Load Balancers (NGINX, HAProxy)** – для распределения нагрузки.
- **Kafka или RabbitMQ** – для обработки сообщений.

---

## Безопасность WebSockets
WebSockets быстрые, но требуют защиты:
- **HTTPS + WSS** – шифруем трафик.
- **Аутентификация (JWT, токены)** – проверяем пользователей.
- **Rate Limiting** – ограничиваем злоумышленников.
- **CORS** – разрешаем запросы только от доверенных доменов.

Пример безопасного соединения с токеном:
```go
http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !validateToken(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	handleConnection(w, r)
})
```

---

## Заключение
WebSockets – это мощный инструмент для реального времени. Они позволяют сделать приложения быстрыми, отзывчивыми и масштабируемыми.
