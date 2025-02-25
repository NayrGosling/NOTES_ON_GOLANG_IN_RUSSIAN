Основы SQL
  Что такое SQL?
    Язык для работы с реляционными базами данных (Structured Query Language).
    Используется для создания, изменения, управления и запросов данных.
  
Типы команд SQL
  
  DDL (Data Definition Language) — определение структуры:
    CREATE — создание таблицы/БД.
    ALTER — изменение структуры.
    DROP — удаление таблицы/БД.
  
  DML (Data Manipulation Language) — работа с данными:
    SELECT — выбор данных.
    INSERT — добавление данных.
    UPDATE — обновление данных.
    DELETE — удаление данных.
  
  DCL (Data Control Language) — управление доступом:
    GRANT — предоставление прав.
    REVOKE — отзыв прав.
  
  TCL (Transaction Control Language) — управление транзакциями:
    COMMIT — подтверждение изменений.
    ROLLBACK — откат изменений.

Структура таблицы
  Столбцы (Columns): имена и типы данных (INT, VARCHAR, DATE и т. д.).
  Строки (Rows): Конкретные записи.
  Ключи:
    Primary Key (PK) — уникальный идентификатор строки.
    Foreign Key (FK) — ссылка на PK другой таблицы.

Основные команды SQL
  1. SELECT — выбор данных
    Синтаксис: SELECT колонка1, колонка2 FROM таблица WHERE условие;
    Примеры:
      Все записи: SELECT * FROM employees;
      Условие: SELECT name, salary FROM employees WHERE salary > 50000;
      Сортировка: SELECT * FROM employees ORDER BY salary DESC;
      Ограничение: SELECT * FROM employees LIMIT 5;

  2. Агрегатные функции
    COUNT() — подсчет строк.
    SUM() — сумма значений.
    AVG() — среднее значение.
    MAX()/MIN() — максимум/минимум.
    Пример: 
      SELECT department, AVG(salary) FROM employees GROUP BY department;

  3. GROUP BY и HAVING
    GROUP BY: группировка данных.
    HAVING: фильтр для групп (аналог WHERE для агрегатов).
    Пример: 
      SELECT department, COUNT(*) FROM employees GROUP BY department HAVING COUNT(*) > 10;

  4. JOIN — объединение таблиц
    Типы:
    INNER JOIN — только совпадающие записи.
    LEFT JOIN — все из левой таблицы + совпадения из правой.
    RIGHT JOIN — все из правой таблицы + совпадения из левой.
    FULL JOIN — все из обеих таблиц.
    CROSS JOIN — все комбинации.
    Пример: 
      SELECT e.name, d.dept_name FROM employees e INNER JOIN departments d ON e.dept_id = d.id;

  5. INSERT — добавление данных
    Синтаксис: INSERT INTO таблица (колонка1, колонка2) VALUES (значение1, значение2);
    Пример: 
      INSERT INTO employees (name, salary) VALUES ('Иван', 60000);

  6. UPDATE — обновление данных
    Синтаксис: UPDATE таблица SET колонка = значение WHERE условие;
    Пример: 
      UPDATE employees SET salary = 70000 WHERE name = 'Иван';

  7. DELETE — удаление данных
    Синтаксис: DELETE FROM таблица WHERE условие;
    Пример: 
      DELETE FROM employees WHERE salary < 30000;

  8. CREATE TABLE — создание таблицы
    Синтаксис:

      CREATE TABLE employees (
          id INT PRIMARY KEY,
          name VARCHAR(50),
          salary INT,
          dept_id INT,
          FOREIGN KEY (dept_id) REFERENCES departments(id)
      );

Транзакции
  BEGIN TRANSACTION — начало.
  COMMIT — сохранить изменения.
  ROLLBACK — откатить изменения.
  Пример:
    BEGIN TRANSACTION;
    UPDATE employees SET salary = salary * 1.1;
    COMMIT;
















