# Использование Grafana в Golang

## Введение

**Grafana** — это мощный инструмент визуализации данных с открытым исходным кодом. Он используется для мониторинга, анализа и построения графиков на основе данных из различных источников, включая Prometheus, InfluxDB, MySQL, PostgreSQL и другие.

В этой лекции мы разберём, как настроить и использовать Grafana для мониторинга метрик в приложении на Go.

---

## 1. Установка и запуск Grafana

### Установка Grafana

#### Установка в Docker:
```sh
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

#### Установка на Windows/macOS/Linux:
1. Перейдите на [официальный сайт Grafana](https://grafana.com/grafana/download).
2. Скачайте и установите соответствующую версию для вашей ОС.
3. Запустите Grafana:
   ```sh
   systemctl start grafana-server
   ```
4. По умолчанию Grafana доступна по адресу: [http://localhost:3000](http://localhost:3000)
   - Логин: `admin`
   - Пароль: `admin` (при первом входе требуется сменить пароль)

---

## 2. Интеграция Grafana с Prometheus

### Добавление источника данных в Grafana
1. Откройте Grafana в браузере: [http://localhost:3000](http://localhost:3000)
2. Перейдите в **Settings → Data Sources**
3. Нажмите **Add data source**
4. Выберите **Prometheus**
5. В поле `URL` укажите адрес сервера Prometheus (по умолчанию: `http://localhost:9090`)
6. Нажмите **Save & Test**

---

## 3. Настройка метрик в Go

### Установка библиотеки Prometheus для Go
```sh
go get github.com/prometheus/client_golang/prometheus
 go get github.com/prometheus/client_golang/prometheus/promhttp
```

### Создание простого API с метриками
```go
package main

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
)

// Создаём метрику - счётчик запросов
var requestCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Общее количество HTTP-запросов",
	},
)

func main() {
	// Регистрируем метрику
	prometheus.MustRegister(requestCount)

	// Обработчик запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestCount.Inc()
		w.Write([]byte("Hello, Grafana!"))
	})

	// Экспорт метрик на /metrics
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Запуск сервера на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Запустите приложение и откройте `http://localhost:8080/metrics`, чтобы проверить, что метрики работают.

---

## 4. Визуализация данных в Grafana

### Создание дашборда
1. В Grafana перейдите в **Dashboards → New Dashboard**
2. Добавьте новую панель (**Add new panel**)
3. В поле `Query` введите запрос PromQL:
   ```promql
   http_requests_total
   ```
4. Настройте визуализацию (график, таблица и т.д.)
5. Сохраните дашборд

Теперь Grafana будет отображать количество HTTP-запросов к вашему сервису!

---

## 5. Настройка алертов (уведомлений)

Grafana позволяет настроить уведомления, если метрики выходят за пределы нормы.

### Создание алерта
1. Откройте панель с метрикой
2. Перейдите во вкладку **Alert**
3. Нажмите **Create Alert Rule**
4. Укажите условие, например:
   ```
   WHEN avg(http_requests_total) > 1000
   ```
5. Выберите канал уведомлений (Email, Slack, Telegram и т.д.)
6. Сохраните правило

---

## 6. Заключение

В этой лекции мы разобрали:
✅ Установку и запуск Grafana
✅ Интеграцию Grafana с Prometheus
✅ Мониторинг метрик из Go-приложения
✅ Настройку дашбордов и графиков
✅ Создание алертов для уведомлений

Теперь ваше приложение можно полноценно мониторить и анализировать!

---

### Полезные ссылки:
- [Официальная документация Grafana](https://grafana.com/docs/)
- [Grafana + Prometheus](https://prometheus.io/docs/visualization/grafana/)

