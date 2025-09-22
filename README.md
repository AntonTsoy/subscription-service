# subscription-service

Go REST-сервис для агрегации данных об онлайн-подписках пользователей. Позволяет создавать, обновлять, удалять и просматривать подписки, а также рассчитывать общую стоимость подписок с дополнительной фильтрацией по пользователю и сервису подписки.

## Используемые технологии

- **Go 1.23** — основной язык
- **Chi** — роутер для HTTP API
- **sqlx** — работа с БД
- **PostgreSQL** — база данных
- **golang-migrate** — миграции
- **Swagger (swaggo)** — документация API
- **Docker + docker-compose** — контейнеризация и оркестрация

## Запуск

1. Клонирование репозитория
```bash
git clone https://github.com/AntonTsoy/subscription-service.git
cd subscription-service
```

2. Заполнить файл .env, **пример моего**:
```bash
# Контейнер PostgreSQL
POSTGRES_USER=service
POSTGRES_PASSWORD=never_guess_password
POSTGRES_DB=subs

# Подключение приложения к БД
DB_HOST=db
DB_PORT=5432
DB_USER=${POSTGRES_USER}
DB_PASSWORD=${POSTGRES_PASSWORD}
DB_NAME=${POSTGRES_DB}
DB_SSL_MODE=disable
```


3. Собрать и запустить сервис
```bash
docker-compose up --build
```

4. Приложение будет доступно на http://localhost:8080

5. Swagger-документация доступна по адресу http://localhost:8080/swagger/index.html

## API

Основные эндпоинты:

- `POST /subscriptions` — создать подписку
- `GET /subscriptions/{id}` — получить подписку по ID
- `GET /subscriptions?limit={limit}&offset={offset}` — получить список подписок (с пагинацией)
- `PUT /subscriptions/{id}` — обновить подписку
- `DELETE /subscriptions/{id}` — удалить подписку
- `GET /subscriptions/{start}/{end}/total-cost?user_id={user_id}&service_name={service_name}` — сумма стоимости подписко за период, с возможной фильтрацией по пользователю и сервису

## Архитектура

- `cmd/app` — точка входа
- `internal/config` — конфигурация
- `internal/database` — подключение к базе
- `internal/repository` — работа с БД
- `internal/service` — бизнес-логика
- `internal/transport/handler` — HTTP-обработчики
- `internal/transport/dto` — DTO для входных/выходных данных
- `internal/transport/logger` - Logger для запросов клиента
- `migrations` — SQL-миграции

## Миграции

Для работы с миграциями используется контейнер в docker compose от `golang-migrate`
