# Работа с SQL-транзакциями в PostgreSQL с использованием Golang

## Введение

Транзакции — это фундаментальная концепция в реляционных базах данных, таких как PostgreSQL. Они позволяют выполнять набор операций как единое целое, обеспечивая свойства ACID (Atomicity, Consistency, Isolation, Durability). В этой лекции мы разберём, что такое транзакции, как они работают в PostgreSQL, и как их использовать в приложении на Go с библиотекой `database/sql`.

Мы рассмотрим:
- Основы транзакций в PostgreSQL.
- Подключение к PostgreSQL из Go.
- Реализацию транзакций с использованием `Begin()`, `Commit()` и `Rollback()`.
- Практические примеры и обработку ошибок.

---

## Что такое транзакции в SQL?

Транзакция — это последовательность операций с базой данных, которая рассматривается как единое целое. Если все операции выполнены успешно, изменения фиксируются (`COMMIT`). Если что-то пошло не так, изменения откатываются (`ROLLBACK`).

### Свойства ACID
1. **Atomicity (Атомарность)**: Все операции либо выполняются полностью, либо не выполняются вовсе.
2. **Consistency (Согласованность)**: После завершения транзакции база данных остаётся в корректном состоянии.
3. **Isolation (Изоляция)**: Транзакции не мешают друг другу до завершения.
4. **Durability (Долговечность)**: Зафиксированные изменения сохраняются даже при сбоях.

### Пример транзакции в PostgreSQL
```sql
BEGIN;
INSERT INTO accounts (user_id, balance) VALUES (1, 1000);
UPDATE accounts SET balance = balance - 200 WHERE user_id = 1;
INSERT INTO transactions (user_id, amount) VALUES (1, -200);
COMMIT;
```

- Если любая из операций завершится ошибкой (например, из-за нарушения ограничений), можно выполнить `ROLLBACK`, и все изменения будут отменены.

---

## Подключение к PostgreSQL в Go

Для работы с PostgreSQL в Go используется стандартная библиотека `database/sql` и драйвер, например, `github.com/lib/pq`. Сначала установим драйвер:

```bash
go get github.com/lib/pq
```

### Пример подключения
```go
package main

import (
    "database/sql"
    _ "github.com/lib/pq" // Регистрация драйвера
    "log"
)

func main() {
    connStr := "user=postgres password=secret dbname=mydb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Successfully connected to PostgreSQL!")
}
```

- `sql.Open()` создаёт объект базы данных (`*sql.DB`), который является пулом соединений.
- `_ "github.com/lib/pq"` импортирует драйвер анонимно, чтобы он зарегистрировался в `database/sql`.

---

## Транзакции в Go

В Go транзакции реализуются через методы `*sql.DB`:
- `Begin()` — начинает новую транзакцию и возвращает объект `*sql.Tx`.
- `Commit()` — фиксирует изменения.
- `Rollback()` — откатывает изменения.

Объект `*sql.Tx` используется для выполнения запросов в рамках транзакции. После вызова `Commit()` или `Rollback()` транзакция завершается, и повторное её использование невозможно.

### Базовый пример транзакции
Предположим, у нас есть таблица `accounts`:
```sql
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    balance DECIMAL NOT NULL
);
```

Теперь реализуем перевод денег между двумя счетами в одной транзакции:

```go
package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    // Начало транзакции
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    // Откладываем Rollback на случай ошибки
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
    }()

    // Снимаем деньги с первого счёта
    _, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE user_id = $2", amount, fromID)
    if err != nil {
        return err
    }

    // Добавляем деньги на второй счёт
    _, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE user_id = $2", amount, toID)
    if err != nil {
        return err
    }

    // Фиксируем транзакцию
    err = tx.Commit()
    if err != nil {
        return err
    }

    log.Println("Transfer completed successfully!")
    return nil
}

func main() {
    connStr := "user=postgres password=secret dbname=mydb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = transferMoney(db, 1, 2, 200.00)
    if err != nil {
        log.Printf("Transfer failed: %v", err)
    }
}
```

#### Разбор:
1. **Начало транзакции**: `db.Begin()` создаёт объект `tx`.
2. **defer tx.Rollback()`**: Если произойдёт ошибка, откат произойдёт автоматически. Мы проверяем `err` перед этим, чтобы не откатывать уже зафиксированную транзакцию.
3. **Выполнение запросов**: Используем `tx.Exec()` для операций в рамках транзакции.
4. **Фиксация**: `tx.Commit()` сохраняет изменения.

---

## Обработка ошибок

Ошибки могут возникнуть на любом этапе: при начале транзакции, выполнении запросов или фиксации. Важно правильно их обрабатывать.

### Пример с проверкой баланса
Добавим проверку, чтобы нельзя было снять больше денег, чем есть на счету:

```go
func transferMoneyWithCheck(db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
    }()

    // Проверяем баланс
    var balance float64
    err = tx.QueryRow("SELECT balance FROM accounts WHERE user_id = $1 FOR UPDATE", fromID).Scan(&balance)
    if err != nil {
        return err
    }
    if balance < amount {
        err = fmt.Errorf("insufficient funds: %v < %v", balance, amount)
        return err
    }

    // Снимаем деньги
    _, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE user_id = $2", amount, fromID)
    if err != nil {
        return err
    }

    // Добавляем деньги
    _, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE user_id = $2", amount, toID)
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        return err
    }

    log.Println("Transfer completed successfully!")
    return nil
}
```

- `FOR UPDATE` блокирует строку, чтобы избежать конкурентных изменений.
- Если баланс недостаточен, возвращаем ошибку, и `defer` вызовет `Rollback()`.

---

## Использование `RETURNING` в транзакциях

Комбинируем транзакции с `RETURNING` для получения данных:

```go
func insertAndGetID(db *sql.DB, userID int, balance float64) (int, error) {
    tx, err := db.Begin()
    if err != nil {
        return 0, err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
    }()

    var id int
    err = tx.QueryRow("INSERT INTO accounts (user_id, balance) VALUES ($1, $2) RETURNING id", userID, balance).Scan(&id)
    if err != nil {
        return 0, err
    }

    err = tx.Commit()
    if err != nil {
        return 0, err
    }

    return id, nil
}
```

- `RETURNING id` возвращает сгенерированный ID, который мы считываем через `Scan()`.

---

## Практические советы

1. **Ограничение времени**: Используйте `context` с таймаутом для транзакций:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   tx, err := db.BeginTx(ctx, nil)
   ```
2. **Изоляция**: Укажите уровень изоляции при необходимости через `sql.TxOptions`:
   ```go
   tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
   ```
3. **Логирование**: Логируйте успешные и неуспешные операции для отладки.

---

## Заключение

Транзакции в PostgreSQL и Go позволяют безопасно и эффективно управлять изменениями в базе данных. Используя `database/sql`, вы можете легко реализовать атомарные операции, обрабатывать ошибки и интегрировать дополнительные возможности, такие как `RETURNING`. Освоив этот подход, вы сможете строить надёжные приложения с корректной обработкой данных.