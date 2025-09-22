# subscription-service

Go REST-сервис для агрегации данных об онлайн-подписках пользователей

## Мой `.env` файл
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
