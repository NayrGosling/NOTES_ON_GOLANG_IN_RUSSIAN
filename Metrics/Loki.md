# Использование Loki в Golang

## Введение

**Loki** — это система агрегирования логов от создателей Grafana, предназначенная для эффективного сбора, хранения и анализа логов. Loki работает по принципу «pull-based», что делает его удобным для использования с микросервисами.

В этом уроке мы разберём:
- Установку и настройку Loki
- Интеграцию Loki с Grafana
- Использование клиента **Promtail** для сбора логов
- Отправку логов из Go-приложения

---

## 1. Установка и настройка Loki

### Установка Loki в Docker
```sh
docker run -d --name=loki -p 3100:3100 grafana/loki:latest
```

### Проверка работы Loki
После запуска Loki должен быть доступен по адресу `http://localhost:3100`. Чтобы проверить статус, выполните:
```sh
curl -s http://localhost:3100/ready
```
Ответ `ready` означает, что Loki работает.

---

## 2. Установка и настройка Promtail

**Promtail** — это агент для сбора логов, который отправляет их в Loki.

### Установка Promtail в Docker
```sh
docker run -d --name=promtail -v $(pwd)/promtail-config.yml:/etc/promtail/config.yml -p 9080:9080 grafana/promtail:latest -config.file=/etc/promtail/config.yml
```

### Настройка `promtail-config.yml`
Создайте файл `promtail-config.yml`:
```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: "varlogs"
    static_configs:
      - targets:
          - localhost
        labels:
          job: "varlogs"
          host: "localhost"
          __path__: /var/log/*.log
```

Запустите Promtail и убедитесь, что он отправляет логи в Loki.

---

## 3. Интеграция Loki с Grafana

### Добавление Loki как источника данных
1. Перейдите в Grafana ([http://localhost:3000](http://localhost:3000))
2. В **Settings → Data Sources** выберите **Loki**
3. В поле `URL` укажите `http://localhost:3100`
4. Нажмите **Save & Test**

Теперь Grafana может запрашивать логи из Loki.

---

## 4. Логирование в Go с Loki

### Установка клиента Loki
```sh
go get github.com/go-kit/kit/log
```

### Пример отправки логов в Loki
```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LokiPayload struct {
	Streams []LokiStream `json:"streams"`
}

func sendLogToLoki(logMessage string) {
	url := "http://localhost:3100/loki/api/v1/push"
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	logEntry := LokiPayload{
		Streams: []LokiStream{
			{
				Stream: map[string]string{
					"job":  "golang_app",
					"level": "info",
				},
				Values: [][]string{{timestamp, logMessage}},
			},
		},
	}

	jsonData, _ := json.Marshal(logEntry)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при отправке логов в Loki:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Лог отправлен в Loki с кодом ответа:", resp.StatusCode)
}

func main() {
	sendLogToLoki("Приложение запущено!")
}
```

Запустите приложение и проверьте, что логи отображаются в Grafana (`Explore → Loki`).

---

## 5. Заключение

В этом уроке мы:
✅ Установили и запустили Loki
✅ Настроили Promtail для сбора логов
✅ Интегрировали Loki с Grafana
✅ Научились отправлять логи из Go-приложения

Теперь ваше приложение может вести централизованный сбор логов! 🚀

---

### Полезные ссылки:
- [Документация Loki](https://grafana.com/docs/loki/latest/)
- [Loki API](https://grafana.com/docs/loki/latest/api/)

