# Использование Prometheus в Golang

## Введение

Prometheus — это система мониторинга и оповещения с открытым исходным кодом, разработанная в SoundCloud. Она предназначена для сбора и хранения метрик в формате временных рядов (time-series). Prometheus хорошо интегрируется с микросервисной архитектурой и поддерживает автоматическое обнаружение сервисов.

В этом уроке мы разберёмся, как интегрировать Prometheus в приложение на Go.

---

## 1. Установка и настройка Prometheus

### Установка Prometheus
1. Перейдите на [официальный сайт](https://prometheus.io/download/).
2. Скачайте соответствующий дистрибутив для вашей ОС.
3. Разархивируйте файлы и запустите `prometheus`:
   ```sh
   ./prometheus --config.file=prometheus.yml
   ```

### Настройка `prometheus.yml`
Создайте файл `prometheus.yml` с таким содержимым:
```yaml
global:
  scrape_interval: 15s # Как часто собирать метрики

scrape_configs:
  - job_name: 'golang_app'
    static_configs:
      - targets: ['localhost:8080']
```

---

## 2. Подключение Prometheus в Go

### Установка библиотеки
В Go есть официальная библиотека для работы с Prometheus:
```sh
 go get github.com/prometheus/client_golang/prometheus
 go get github.com/prometheus/client_golang/prometheus/promhttp
```

### Добавление метрик
Создадим простое веб-приложение, экспортирующее метрики:
```go
package main

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
)

// Создание метрики
var requestCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Общее количество HTTP-запросов",
	},
)

func main() {
	// Регистрируем метрику в Prometheus
	prometheus.MustRegister(requestCount)

	// Обработчик запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestCount.Inc() // Увеличиваем счетчик при каждом запросе
		w.Write([]byte("Hello, Prometheus!"))
	})

	// Регистрируем обработчик метрик
	http.Handle("/metrics", promhttp.Handler())

	// Запускаем сервер
	log.Println("Запуск сервера на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Запустите приложение и откройте `http://localhost:8080/metrics` в браузере.

---

## 3. Запросы в Prometheus

После того как Prometheus начал собирать метрики, откройте веб-интерфейс `http://localhost:9090` и выполните запрос:
```promql
http_requests_total
```
Это покажет общее количество запросов.

---

## 4. Визуализация в Grafana

Prometheus хорошо интегрируется с [Grafana](https://grafana.com/):
1. Установите Grafana.
2. Добавьте Prometheus как источник данных (`http://localhost:9090`).
3. Создайте новый дашборд и используйте запрос `http_requests_total`.

---

## Заключение
Мы разобрали основы интеграции Prometheus в Go:
✅ Установили Prometheus
✅ Настроили `prometheus.yml`
✅ Подключили клиентскую библиотеку в Go
✅ Создали и зарегистрировали метрики
✅ Настроили визуализацию в Grafana

Дальше можно добавить кастомные метрики, алерты и масштабируемый сбор данных!

---

### Полезные ссылки:
- [Официальная документация Prometheus](https://prometheus.io/docs/)
- [Prometheus в Go](https://github.com/prometheus/client_golang)
- [Grafana](https://grafana.com/)

