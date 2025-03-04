# Основы SQL на примере PostgreSQL

## Введение

SQL (Structured Query Language) — это язык запросов, используемый для работы с реляционными базами данных. PostgreSQL — одна из самых мощных и популярных open-source СУБД, которая полностью поддерживает SQL. В этой лекции мы разберём основы SQL: создание таблиц, добавление данных, выборку, обновление и удаление, а также более сложные конструкции, такие как объединения (JOIN), подзапросы и агрегации.

Цели лекции:
- Понять структуру базовых SQL-запросов.
- Освоить работу с данными в PostgreSQL.
- Изучить сложные примеры для практического применения.

---

## 1. Создание таблиц (`CREATE TABLE`)

SQL начинается с определения структуры данных. Таблицы — это основа реляционной базы данных.

### Пример: Создание таблицы пользователей
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

- `id SERIAL PRIMARY KEY`: Автоинкрементное поле, уникальный идентификатор.
- `name TEXT NOT NULL`: Имя пользователя, не может быть пустым.
- `email TEXT UNIQUE`: Уникальный email.
- `created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`: Время создания строки, по умолчанию текущая дата и время.

#### Проверка структуры
```sql
\d users
```

---

## 2. Добавление данных (`INSERT`)

Теперь добавим данные в таблицу.

### Простой пример
```sql
INSERT INTO users (name, email)
VALUES ('Иван Иванов', 'ivan@example.com');
```

Результат: Добавлена одна строка, `id` автоматически сгенерирован (например, 1), `created_at` установлено текущее время.

### Добавление нескольких строк
```sql
INSERT INTO users (name, email)
VALUES 
    ('Мария Петрова', 'maria@example.com'),
    ('Алексей Сидоров', 'alexey@example.com');
```

### Использование `RETURNING`
```sql
INSERT INTO users (name, email)
VALUES ('Ольга Смирнова', 'olga@example.com')
RETURNING id, name;
```

**Результат:**
```
 id |    name
----+------------
  4 | Ольга Смирнова
```

---

## 3. Выборка данных (`SELECT`)

`SELECT` используется для получения данных из таблицы.

### Простая выборка
```sql
SELECT * FROM users;
```

**Результат:**
```
 id |    name         |        email          |      created_at
----+-----------------+----------------------+------------------------
  1 | Иван Иванов     | ivan@example.com     | 2025-03-03 10:00:00
  2 | Мария Петрова   | maria@example.com    | 2025-03-03 10:01:00
  3 | Алексей Сидоров | alexey@example.com   | 2025-03-03 10:02:00
  4 | Ольга Смирнова  | olga@example.com     | 2025-03-03 10:03:00
```

- `*` возвращает все столбцы.

### Выборка с условием (`WHERE`)
```sql
SELECT name, email 
FROM users 
WHERE id > 2;
```

**Результат:**
```
    name         |        email
-----------------+-------------------
 Алексей Сидоров | alexey@example.com
 Ольга Смирнова  | olga@example.com
```

---

## 4. Обновление данных (`UPDATE`)

Модифицируем существующие записи.

### Простое обновление
```sql
UPDATE users
SET email = 'ivan_new@example.com'
WHERE id = 1;
```

- Изменён email для пользователя с `id = 1`.

### Обновление с возвратом
```sql
UPDATE users
SET name = UPPER(name)
WHERE id = 2
RETURNING *;
```

**Результат:**
```
 id |    name       |        email       |      created_at
----+---------------+--------------------+------------------------
  2 | МАРИЯ ПЕТРОВА | maria@example.com  | 2025-03-03 10:01:00
```

- `UPPER(name)` преобразует имя в верхний регистр.

---

## 5. Удаление данных (`DELETE`)

Удаляем строки из таблицы.

### Простое удаление
```sql
DELETE FROM users
WHERE id = 3;
```

- Удалён пользователь с `id = 3`.

### Удаление с возвратом
```sql
DELETE FROM users
WHERE email = 'olga@example.com'
RETURNING name, email;
```

**Результат:**
```
    name        |        email
----------------+-------------------
 Ольга Смирнова | olga@example.com
```

---

## 6. Сложные примеры

Теперь перейдём к более сложным конструкциям.

### Создание второй таблицы: заказы
```sql
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    amount DECIMAL NOT NULL,
    order_date DATE DEFAULT CURRENT_DATE
);

INSERT INTO orders (user_id, amount)
VALUES 
    (1, 150.50),
    (2, 200.00),
    (1, 75.25);
```

### 6.1. Объединение таблиц (`JOIN`)
Получим имена пользователей и их заказы:
```sql
SELECT u.name, o.order_id, o.amount
FROM users u
JOIN orders o ON u.id = o.user_id;
```

**Результат:**
```
    name       | order_id | amount
---------------+----------+--------
 Иван Иванов   |    1     | 150.50
 Мария Петрова |    2     | 200.00
 Иван Иванов   |    3     | 75.25
```

- `JOIN` связывает таблицы по `id` и `user_id`.

#### Использование `LEFT JOIN`
Если нужны все пользователи, даже без заказов:
```sql
SELECT u.name, o.order_id, o.amount
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;
```

---

### 6.2. Агрегация (`GROUP BY`, `HAVING`)
Посчитаем общую сумму заказов по пользователям:
```sql
SELECT u.name, COUNT(o.order_id) AS order_count, SUM(o.amount) AS total_amount
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.name
HAVING SUM(o.amount) > 100;
```

**Результат:**
```
    name       | order_count | total_amount
---------------+-------------+--------------
 Иван Иванов   |     2       |   225.75
 Мария Петрова |     1       |   200.00
```

- `COUNT` считает количество заказов.
- `SUM` суммирует суммы.
- `HAVING` фильтрует группы с общей суммой больше 100.

---

### 6.3. Подзапросы
Найдём пользователей, у которых есть заказы:
```sql
SELECT name, email
FROM users
WHERE id IN (SELECT user_id FROM orders);
```

**Результат:**
```
    name       |        email
---------------+--------------------
 Иван Иванов   | ivan_new@example.com
 Мария Петрова | maria@example.com
```

- Подзапрос возвращает `user_id` из `orders`, а внешний запрос фильтрует пользователей.

---

### 6.4. Условные выражения (`CASE`)
Добавим категорию суммы заказа:
```sql
SELECT 
    u.name, 
    o.amount,
    CASE 
        WHEN o.amount > 150 THEN 'Большой заказ'
        WHEN o.amount > 50 THEN 'Средний заказ'
        ELSE 'Маленький заказ'
    END AS order_category
FROM users u
JOIN orders o ON u.id = o.user_id;
```

**Результат:**
```
    name       | amount | order_category
---------------+--------+----------------
 Иван Иванов   | 150.50 | Большой заказ
 Мария Петрова | 200.00 | Большой заказ
 Иван Иванов   | 75.25  | Средний заказ
```

---

## 7. Полезные советы

1. **Индексы**: Создавайте индексы для часто используемых столбцов (например, `CREATE INDEX ON orders(user_id);`) для ускорения запросов.
2. **Ограничения**: Используйте `LIMIT` и `OFFSET` для пагинации:
   ```sql
   SELECT * FROM users LIMIT 2 OFFSET 1;
   ```
3. **Безопасность**: При работе с параметрами в приложении используйте placeholders (`$1`, `$2`) для защиты от SQL-инъекций.

---

## Заключение

SQL в PostgreSQL предоставляет мощные инструменты для управления данными: от простых операций вставки и выборки до сложных аналитических запросов с объединениями и агрегациями. Освоив базовые команды (`SELECT`, `INSERT`, `UPDATE`, `DELETE`) и продвинутые конструкции (`JOIN`, `GROUP BY`, подзапросы), вы сможете эффективно работать с данными в реальных проектах.
