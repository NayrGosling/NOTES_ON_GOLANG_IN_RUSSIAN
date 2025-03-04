# Логирование с использованием Promtail, Loki и Golang

## Введение

**Promtail** — это агент логирования, который собирает логи с файловой системы и отправляет их в **Loki**, хранилище логов от Grafana. Вместе с **Golang** они позволяют эффективно управлять логами приложений.

В этой лекции мы разберем:
- Основы Loki и Promtail
- Установку и настройку Loki и Promtail
- Интеграцию с Golang
- Практическое применение и примеры

---

## 1. Основы Loki и Promtail

### 🔹 **Что такое Loki?**
Loki — это система логирования от **Grafana**, работающая аналогично Prometheus, но для логов:
- Эффективно индексирует логи (только метаданные).
- Поддерживает структурированные и неструктурированные логи.
- Интегрируется с Grafana для визуализации.

### 🔹 **Что такое Promtail?**
Promtail — это агент, который:
- Читает логи из файлов и контейнеров.
- Добавляет метаданные (например, `job`, `instance`).
- Отправляет логи в Loki через HTTP API.

---

## 2. Установка и настройка Loki и Promtail

### 📌 **Установка Loki**
Скачиваем и запускаем Loki:
```sh
wget https://github.com/grafana/loki/releases/latest/download/loki-linux-amd64.zip
unzip loki-linux-amd64.zip
chmod +x loki-linux-amd64
./loki-linux-amd64 -config.file=loki-config.yaml
```

Конфигурация `loki-config.yaml`:
```yaml
auth_enabled: false
server:
  http_listen_port: 3100
ingester:
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h
```

### 📌 **Установка Promtail**
Скачиваем и запускаем Promtail:
```sh
wget https://github.com/grafana/loki/releases/latest/download/promtail-linux-amd64.zip
unzip promtail-linux-amd64.zip
chmod +x promtail-linux-amd64
./promtail-linux-amd64 -config.file=promtail-config.yaml
```

Конфигурация `promtail-config.yaml`:
```yaml
server:
  http_listen_port: 9080
positions:
  filename: /tmp/positions.yaml
clients:
  - url: http://localhost:3100/loki/api/v1/push
scrape_configs:
  - job_name: "golang-logs"
    static_configs:
      - targets:
          - localhost
        labels:
          job: "golang"
          __path__: "/var/log/*.log"
```

---

## 3. Интеграция с Golang

### 📌 **Настройка логирования в Golang**

Устанавливаем клиент для Loki:
```sh
go get github.com/grafana/loki-client-go/lokiclient
```

Пример кода для логирования в Loki:
```go
package main

import (
	"log"
	"os"
	"github.com/grafana/loki-client-go/lokiclient"
	"github.com/grafana/loki-client-go/model"
	"time"
)

func main() {
	logger := lokiclient.New(lokiclient.Config{
		URL: "http://localhost:3100/loki/api/v1/push",
	})

	entry := model.Entry{
		Timestamp: time.Now(),
		Line:      "Привет, Loki!",
	}

	err := logger.Handle(entry)
	if err != nil {
		log.Fatalf("Ошибка логирования: %v", err)
	}

	log.Println("Лог отправлен в Loki")
}
```

---

## 4. Визуализация логов в Grafana

1️⃣ Установите Grafana и запустите ее:
```sh
docker run -d -p 3000:3000 grafana/grafana
```
2️⃣ Откройте `http://localhost:3000`, войдите (`admin/admin`).
3️⃣ Добавьте Loki как источник данных (`Settings -> Data Sources -> Loki`).
4️⃣ Создайте дашборд и выполните запрос:
```logql
{job="golang"} |= "Привет"
```

---

## Заключение

✅ **Promtail** собирает логи и отправляет их в **Loki**.
✅ **Loki** хранит логи и позволяет быстро их искать.
✅ **Golang** может логировать напрямую в **Loki**.
✅ **Grafana** визуализирует логи и упрощает мониторинг.

Сочетание Promtail, Loki и Golang делает логирование мощным и удобным инструментом! 🚀

