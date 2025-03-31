-- Создание таблицы пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

-- Создание таблицы списков задач
CREATE TABLE todo_lists (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255)
);

-- Промежуточная таблица "связка пользователей и списков"
CREATE TABLE users_lists (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    list_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (list_id) REFERENCES todo_lists(id) ON DELETE CASCADE,
    UNIQUE (user_id, list_id)
);

-- Создание таблицы задач
CREATE TABLE todo_items(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    done BOOLEAN NOT NULL DEFAULT FALSE
);

-- Промежуточная таблица "связка списков и задач"
CREATE TABLE lists_items (
    id SERIAL PRIMARY KEY,
    list_id INT NOT NULL,
    item_id INT NOT NULL,
    FOREIGN KEY (list_id) REFERENCES todo_lists(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES todo_items(id) ON DELETE CASCADE,
    UNIQUE (list_id, item_id)
);