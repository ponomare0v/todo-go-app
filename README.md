# Todo Go App

## Описание
REST API для управления списками задач, написанный на Go с использованием PostgreSQL.

## Запуск с помощью Docker

### 1. Клонирование репозитория
```sh
git clone https://github.com/ponomare0v/todo-go-app.git
cd todo-go-app
```

### 2. Создание `.env`
```sh
cp .env.example .env
```

### 3. Запуск контейнеров
```sh
docker-compose up --build
```

### 4. Проверка работы
После запуска будет доступна Swagger документация по адресу:
```
http://localhost:8000/swagger/index.html#/
```