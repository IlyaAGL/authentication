# Authentication Service

Сервис аутентификации на Go с использованием JWT, PostgreSQL и Gin.

## Запуск проекта

### 1. Клонируй репозиторий

```bash
git clone https://github.com/IlyaAGL/authentication.git
cd authentication
```

### 2. Подними PostgreSQL и pgAdmin через Docker

```bash
cd docker/db
docker-compose up -d
```

- PostgreSQL будет доступен на порту `5432`
- pgAdmin — на `http://localhost:8080`
## ️ Таблица `tokens` в PostgreSQL

```sql
CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    tokenPairID VARCHAR(255) NOT NULL
);
```

## Запуск приложения

```bash
cd authentication/cmd
go run main.go
```

Сервис стартует на `http://localhost:6060`.

##  API

###  Получить Access Token

```http
GET /tokens/:id
```

- **Описание**: генерирует `access token` и записывает зашифрованный refresh-токен в БД.
- **Ответ**:
  - Успех: `200 OK` - возвращает `access token` и записывает в cookie
  - Ошибка: `400 Bad Request`

###  Обновить Access Token

```http
GET /tokens
```

- **Описание**: получает `access_token` из куки, валидирует и создает новый.
- **Ответ**:
  - Успех: `200 OK` - возвращает `access token` и записывает в cookie
  - Ошибка: `401 Unauthorized` или `400 Bad Request`
