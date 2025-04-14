# Algobot

## Описание
**Алгобот** — это Telegram-бот, предназначенный для помощи в ведении учительной деятельности Алгоритмика

Репозиторий состоит из двух приложений:
- `algobot` — основной Telegram-бот
- `migrator` — утилита для применения миграций с помощью [Goose](https://github.com/pressly/goose)

Бот использует:
- [gpt4free-grpc-gateway](https://github.com/LZTD1/gpt4free-grpc-gateway) — gRPC-сервер для обращения к нейросетям
- [telebot-context-router](https://github.com/LZTD1/telebot-context-router) — роутер для удобного управления обработчиками сообщений в Telebot

---

## Структура проекта
```
cmd/
├── algobot/     # Точка входа для Telegram-бота
│   └── main.go
├── migrator/    # Миграции SQLite базы через Goose
│   └── main.go

config/          # Конфигурация
migrations/      # SQL миграции
storage/         # SQLite база данных
protos/          # gRPC proto-файлы
```

---

## Конфигурация

Файл `config.yaml`:
```yaml
env: local # или prod
storage_path: "./storage/storage.db" # путь то файла базы данных
# telegram_token можно указать в .env или через переменные окружения TELEGRAM_TOKEN
migrations_path: "./migrations" # путь до папки с миграциями
grpc: # настройки grpc
  host: "localhost" # адрес grpc сервера с нейронкой 
  port: 50051 # порт grpc сервера
  timeout: 300s # таймаут ответа
rate_limit: # rate limit для запросов в бота
  fill_period: 800ms # период пополнение бакета
  bucket_limit: 6 # величина бакета
backoffice: # конфигурация для работы с бэкофисом
  message_timer: 5m # период опроса новых сообщений от учеников
  retries: 3 # количество ретраев при ошибке от сервера
  retries_timeout: 5s # ожидание между ретраями
  response_timeout: 15s # ожидание ответа от сервера
```

Можно передать путь до файла конфига через аргумент `--config` или переменную окружения `CONFIG_PATH`.

---

## Запуск

### 1. Сборка
```bash
go build -o migrator ./cmd/migrator
go build -o algobot ./cmd/algobot
```

### 2. Применение миграций
```bash
./migrator -migrations-path=./migrations -storage-path=./storage/storage.db
```

### 3. Запуск бота
```bash
./algobot -config ./config/config.yaml
```


---

## Makefile команды
Для удобства сборки и генерации предусмотрен Makefile:

| Цель             | Описание                                                        |
|------------------|-----------------------------------------------------------------|
| `make gen`       | Генерация Go и gRPC кода из `.proto` файлов (`./protos`)        |
| `make grpc-gen`  | То же самое, альтернатива `make gen`                            |
| `make dev`       | Запуск бота в dev-режиме с конфигом `./config/local.yaml`       |
| `make mock-gen`  | Генерация mock'ов в папке `test/`                               |
| `make migrate`   | Запуск миграций через `cmd/migrator`                            |

Пример:
```bash
make migrate
make dev
```

---

## Переменные окружения
| Название          | Описание                      |
|-------------------|-------------------------------|
| `TELEGRAM_TOKEN`  | Токен Telegram-бота           |
| `CONFIG_PATH`     | Путь к файлу конфигурации     |
