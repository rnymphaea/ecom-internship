# ToDo API server

[![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-24-2496ED?logo=docker)](https://www.docker.com/)
---

## Содержание
- [Обзор проекта](#обзор-проекта)
- [Основные возможности](#основные-возможности)
- [Технологический стек](#технологический-стек)
- [Структура проекта](#структура-проекта)
- [API Endpoints](#api-endpoints)
- [Быстрый старт](#быстрый-старт)

---
## Обзор проекта

**Todo API server** — это REST API сервис для управления задачами, написанный на Go, используя возможности стандартной библиотеки языка.

---
## Технологический стек
- **Go 1.25.5** - язык программирования
- **Стандартная библиотека Go** - отсутствие внешних зависимостей
- **Docker** - контейнеризация
- **golangci-lint** - линтинг и проверка кода
---

## Структура проекта

```
.
├── cmd/
│   └── main.go                    # Точка входа в приложение
├── internal/
│   ├── app/                       # Инициализация приложения
│   │   ├── app.go                 # Запуск и graceful shutdown
│   │   └── setup.go               # Настройка зависимостей
│   ├── config/                    # Конфигурация
│   │   ├── config.go              # Загрузка конфигурации
│   │   └── config_test.go         # Тесты конфигурации
│   ├── database/                  # Слой данных
│   │   ├── database.go            # Интерфейс БД
│   │   └── mem/                   # In-memory реализация
│   │       ├── mem.go             # Структура хранилища
│   │       ├── todo.go            # CRUD операции
│   │       └── todo_test.go       # Тесты хранилища
│   ├── logger/                    # Логирование
│   │   ├── logger.go              # Интерфейс логгера
│   │   └── std/                   # Реализация с стандартной библиотекой
│   │       └── logger.go          
│   ├── model/                     # Модели данных
│   │   └── model.go               
│   ├── pkg/                       # Вспомогательные пакеты
│   │   └── httputils/             # HTTP утилиты
│   │       └── utils.go           # Работа с контекстом
│   └── server/                    # HTTP сервер
│       ├── handler/               # Обработчики запросов
│       │   ├── handler.go         # Основные обработчики
│       │   └── handler_test.go    # Тесты обработчиков
│       ├── middleware.go          
│       ├── router.go              # Маршрутизация
│       └── server.go              # HTTP сервер
├── .dockerignore                  
├── .gitignore                     
├── .golangci.yaml                 
├── Dockerfile                     
├── go.mod                         
└── Makefile                       # Утилиты сборки и тестирования
```
---

## API Endpoints

### `GET /todos`
Получить список всех задач.

**Ответ:**
```json
{
  "todos": [
    {
      "id": 1,
      "caption": "Купить продукты",
      "description": "Молоко, хлеб, яйца",
      "is_completed": false,
      "created_at": "2025-12-29T10:30:00Z",
      "updated_at": "2025-12-29T10:30:00Z"
    }
  ]
}
```
---
### `GET /todos/{id}`
Получить задачу по ID.

**Ответ:** `200 OK`
```json
{
  "id": 1,
  "caption": "Купить продукты",
  "description": "Молоко, хлеб, яйца",
  "is_completed": false,
  "created_at": "2025-12-29T10:30:00Z",
  "updated_at": "2025-12-29T10:30:00Z"
}
```

**Ошибки:** `404 Not Found` если задача не существует

---

### `POST /todos`
Создать новую задачу.

**Тело запроса:**
```json
{
  "caption": "Новая задача",
  "description": "Описание задачи",
  "is_completed": false
}
```

**Ответ:** `201 Created` с заголовком `Location: host:/todos/{id}`

**Валидация:**
- `caption` не должен быть пустым
- `id` не должен дублироваться

**Ошибки:**
- `400 Bad Request` если `caption` пустой
- `409 Conflict` если `id` уже существует

---

### `PUT /todos/{id}`
Обновить существующую задачу.

**Тело запроса:**
```json
{
  "caption": "Обновленный заголовок",
  "description": "Обновленное описание",
  "is_completed": true
}
```

**Ответ:** `204 No Content`

**Валидация:**
- `caption` не должен быть пустым

**Ошибки:**
- `400 Bad Request` если `caption` пустой
- `404 Not Found` если задача не существует

---

### `DELETE /todos/{id}`
Удалить задачу по ID.

**Ответ:** `204 No Content`

**Ошибки:** `404 Not Found` если задача не существует

## Быстрый старт

### Требования
- Go 1.25.5+
- Docker
- make (опционально)

#### Запуск
```bash
# Клонировать репозиторий
git clone https://github.com/rnymphaea/ecom-internship.git
cd ecom-internship

# Создать файл .env (см. .env.example)

# Запуск с make
make

# Если make нет
docker build -t todo-api . && docker run -p 8080:8080 --env-file .env todo-api
```
#### Использование make 
```bash
make          # или make all - сборка и запуск приложения
make build    # сборка Docker образа
make run      # запуск Docker контейнера
make test     # запуск всех тестов
make lint     # проверка кода линтером
make prepare  # запуск тестов и линтера
make api-test # запуск интеграционных тестов API
```
