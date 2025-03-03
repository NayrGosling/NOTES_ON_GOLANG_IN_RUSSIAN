# ACID & BASE

## 1. Введение
В мире баз данных два важных набора свойств определяют их поведение: **ACID** и **BASE**.
- **ACID** (Atomicity, Consistency, Isolation, Durability) гарантирует надежность транзакций.
- **BASE** (Basically Available, Soft state, Eventual consistency) делает упор на масштабируемость и доступность.

Эти концепции важны при проектировании высоконагруженных систем и распределённых баз данных.

---

## 2. ACID: Гарантия целостности данных
**ACID** — это стандарт для реляционных баз данных (SQL), обеспечивающий корректность транзакций.

### 2.1. Компоненты ACID
| Свойство  | Описание |
|-----------|----------|
| **Atomicity (Атомарность)** | Транзакция либо выполняется полностью, либо не выполняется вовсе. |
| **Consistency (Согласованность)** | Данные переходят из одного корректного состояния в другое. |
| **Isolation (Изоляция)** | Параллельные транзакции не мешают друг другу. |
| **Durability (Долговечность)** | Данные сохраняются навсегда после коммита. |

### 2.2. Пример ACID-транзакции
```sql
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;
```
Если одна из операций не удалась, `ROLLBACK;` отменит все изменения.

### 2.3. Где используется?
- **Реляционные базы данных** (PostgreSQL, MySQL, Oracle)
- **Финансовые системы** (банковские переводы)
- **Заказы и инвентаризация** (e-commerce)

---

## 3. BASE: Масштабируемость в распределённых системах
Когда невозможно строго соблюдать ACID (например, в глобально распределённых базах), используется модель **BASE**.

### 3.1. Компоненты BASE
| Свойство  | Описание |
|-----------|----------|
| **Basically Available (Базовая доступность)** | Система всегда отвечает, пусть даже с устаревшими данными. |
| **Soft state (Мягкое состояние)** | Состояние системы может изменяться со временем без вмешательства пользователя. |
| **Eventual consistency (Конечная согласованность)** | Данные в конечном итоге становятся согласованными. |

### 3.2. Пример BASE-модели
В системе, основанной на BASE (например, **Cassandra**), данные могут быть не мгновенно консистентными:
```sql
INSERT INTO users (id, name) VALUES (1, 'Alice');
-- Данные распространятся по кластерам через некоторое время.
```

### 3.3. Где используется?
- **NoSQL базы данных** (Cassandra, DynamoDB, MongoDB)
- **Социальные сети** (посты и лайки могут обновляться с задержкой)
- **Системы кеширования** (Redis, Memcached)

---

## 4. ACID vs BASE: Выбор подходящей модели
| Характеристика | ACID | BASE |
|---------------|------|------|
| **Целостность данных** | Высокая | Низкая (но со временем приходит в норму) |
| **Доступность** | Может быть снижена ради согласованности | Всегда высокая |
| **Применение** | Банковские системы, платежи, заказы | Большие распределённые системы, соцсети |

### 4.1. CAP-теорема
Согласно **CAP-теореме**, распределённые системы не могут одновременно обеспечивать три свойства:
1. **Consistency (Согласованность)**
2. **Availability (Доступность)**
3. **Partition tolerance (Устойчивость к разделению сети)**

ACID-системы жертвуют доступностью, а BASE-системы жертвуют строгой согласованностью.

---

## 5. Заключение
- **ACID** подходит для систем, где критически важны транзакции (банки, учёт, CRM).
- **BASE** хорош для масштабируемых и высокодоступных решений (NoSQL, Big Data, соцсети).
- В реальных системах часто комбинируют оба подхода, используя **гибридные базы** (например, NewSQL: Google Spanner, CockroachDB).

### Какой подход выбрать?
✅ **Если важна строгая согласованность** → **ACID** (PostgreSQL, MySQL)
✅ **Если важна доступность и масштабируемость** → **BASE** (MongoDB, Cassandra)
✅ **Если нужен баланс** → **NewSQL** (Spanner, TiDB)

