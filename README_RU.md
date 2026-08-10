# Wallets API

[English](README.md) | Русский

Wallets API — учебный REST API на Go для управления пользователями, кошельками и денежными переводами.

Проект построен с разделением ответственности между HTTP-обработчиками, сервисным слоем и слоем работы с базой данных.

API поддерживает регистрацию и авторизацию пользователей, создание нескольких кошельков для одного пользователя, переводы между кошельками, разграничение доступа по ролям, административные операции, SQL-миграции, структурированное логирование, unit-тесты и Swagger/OpenAPI документацию.

## Возможности

### Пользователи

- Регистрация новых пользователей
- Авторизация пользователей
- Хеширование паролей с помощью bcrypt
- JWT-аутентификация
- Роли `user` и `admin`
- Получение списка всех пользователей администратором
- Получение пользователя по ID администратором

### Кошельки

- Создание кошелька
- Создание нескольких кошельков для одного пользователя
- Получение всех кошельков текущего пользователя
- Получение кошелька по ID
- Изменение названия кошелька
- Удаление кошелька
- Привязка кошельков к конкретному пользователю
- Ограничение доступа пользователя только к своим кошелькам
- Получение всех кошельков администратором
- Получение кошельков конкретного пользователя администратором
- Получение любого кошелька по ID администратором

### Переводы

- Создание переводов между кошельками
- Получение истории переводов текущего пользователя
- Получение перевода по ID
- Проверка достаточности баланса
- Запрет перевода между одним и тем же кошельком
- Проверка совместимости валют кошельков
- Выполнение перевода внутри SQL-транзакции
- Получение всех переводов администратором
- Получение любого перевода по ID администратором

## Архитектура

Проект разделён на несколько слоёв:

- `handlers` — обработка HTTP-запросов и формирование HTTP-ответов
- `service` — бизнес-логика приложения
- `repository` — взаимодействие с PostgreSQL
- `models` — структуры данных приложения
- `middleware` — аутентификация, авторизация, логирование и восстановление после panic
- `router` — регистрация HTTP-маршрутов
- `config` — загрузка конфигурации приложения
- `database` — подключение к PostgreSQL
- `response` — формирование единообразных HTTP-ответов
- `auth` — работа с JWT
- `migrations` — SQL-миграции базы данных

Основной поток HTTP-запроса:

```text
Client
  ↓
Router
  ↓
Middleware
  ↓
Handler
  ↓
Service
  ↓
Repository
  ↓
PostgreSQL
```

Такое разделение позволяет отделить HTTP-логику, бизнес-логику и работу с базой данных друг от друга.

## Технологии

- Go
- `net/http`
- PostgreSQL
- SQL
- Docker
- Docker Compose
- JWT
- bcrypt
- `shopspring/decimal`
- `context.Context`
- `log/slog`
- golang-migrate
- Swagger / OpenAPI
- Swaggo
- Git
- GitHub

## Аутентификация

После успешной авторизации пользователь получает JWT-токен.

Для доступа к защищённым endpoint токен передаётся в HTTP-заголовке:

```text
Authorization: Bearer <token>
```

`AuthMiddleware` проверяет JWT, извлекает claims и сохраняет информацию о пользователе в `context.Context` текущего HTTP-запроса.

Handler затем получает данные авторизованного пользователя из контекста.

Для административных маршрутов дополнительно используется `AdminMiddleware`, который проверяет роль пользователя.

## Роли

В приложении используются две роли.

### user

Обычный пользователь может:

- создавать свои кошельки
- получать свои кошельки
- получать свой кошелёк по ID
- изменять свои кошельки
- удалять свои кошельки
- создавать переводы
- просматривать свои переводы

### admin

Администратор дополнительно может:

- получать список всех пользователей
- получать пользователя по ID
- получать список всех кошельков
- получать кошельки конкретного пользователя
- получать любой кошелёк по ID
- получать список всех переводов
- получать любой перевод по ID

## API Endpoints

### Public

| Method | Endpoint | Описание |
|--------|----------|----------|
| POST | `/register` | Регистрация пользователя |
| POST | `/login` | Авторизация и получение JWT |
| GET | `/health` | Проверка работы API |

### Wallets

Требуется JWT-аутентификация.

| Method | Endpoint | Описание |
|--------|----------|----------|
| POST | `/wallets` | Создать кошелёк |
| GET | `/wallets` | Получить свои кошельки |
| GET | `/wallets/{id}` | Получить свой кошелёк по ID |
| PATCH | `/wallets/{id}` | Изменить название кошелька |
| DELETE | `/wallets/{id}` | Удалить кошелёк |

### Transfers

Требуется JWT-аутентификация.

| Method | Endpoint | Описание |
|--------|----------|----------|
| POST | `/transfers` | Создать перевод |
| GET | `/transfers` | Получить свои переводы |
| GET | `/transfers/{id}` | Получить перевод по ID |

### Admin

Требуется JWT-аутентификация и роль `admin`.

| Method | Endpoint | Описание |
|--------|----------|----------|
| GET | `/admin/users` | Получить всех пользователей |
| GET | `/admin/users/{id}` | Получить пользователя по ID |
| GET | `/admin/wallets` | Получить все кошельки |
| GET | `/admin/wallets/{id}` | Получить кошелёк по ID |
| GET | `/admin/users/{id}/wallets` | Получить кошельки конкретного пользователя |
| GET | `/admin/transfers` | Получить все переводы |
| GET | `/admin/transfers/{id}` | Получить перевод по ID |

## Swagger / OpenAPI

Проект содержит Swagger-документацию API.

После запуска приложения Swagger UI доступен по адресу:

```text
http://localhost:8080/swagger/index.html
```

Swagger отображает доступные endpoint, параметры запросов, модели данных и возможные HTTP-ответы.

Для работы с защищёнными endpoint можно использовать кнопку `Authorize` и передать JWT через заголовок `Authorization`.

Генерация Swagger-документации:

```bash
swag init -g ./cmd/api/main.go
```

После генерации создаются файлы:

```text
docs/
├── docs.go
├── swagger.json
└── swagger.yaml
```

## База данных

В проекте используется PostgreSQL.

Основные таблицы:

### users

Хранит пользователей приложения:

- ID
- email
- хеш пароля
- роль
- дата создания
- дата обновления

### wallets

Хранит кошельки пользователей.

Каждый кошелёк связан с конкретным пользователем.

Основные данные:

- ID
- ID пользователя
- название
- валюта
- баланс
- дата создания
- дата обновления

### transfers

Хранит информацию о переводах между кошельками:

- ID
- ID кошелька отправителя
- ID кошелька получателя
- сумма
- статус
- дата создания

## SQL-транзакции

Перевод средств выполняется внутри SQL-транзакции.

Изменение баланса кошелька отправителя, изменение баланса кошелька получателя и создание записи о переводе выполняются как единая операция.

При возникновении ошибки транзакция откатывается.

Это позволяет избежать ситуации, когда деньги были списаны с одного кошелька, но не были зачислены на другой.

## Работа с денежными значениями

Для денежных значений используется пакет:

```text
github.com/shopspring/decimal
```

Использование `decimal.Decimal` позволяет работать с денежными суммами без проблем точности, характерных для вычислений с `float32` и `float64`.

## Миграции

Изменения структуры базы данных управляются через SQL-миграции.

Файлы миграций находятся в директории:

```text
migrations/
```

Для каждой миграции используются два файла:

```text
*.up.sql
*.down.sql
```

`up` применяет изменение базы данных.

`down` откатывает его.

Применить все миграции:

```bash
migrate -path ./migrations -database "postgres://USER:PASSWORD@localhost:5433/wallets_db?sslmode=disable" up
```

Откатить последнюю миграцию:

```bash
migrate -path ./migrations -database "postgres://USER:PASSWORD@localhost:5433/wallets_db?sslmode=disable" down 1
```

Посмотреть текущую версию:

```bash
migrate -path ./migrations -database "postgres://USER:PASSWORD@localhost:5433/wallets_db?sslmode=disable" version
```

## Middleware

В приложении используются несколько middleware.

### AuthMiddleware

Отвечает за JWT-аутентификацию:

- получает JWT из заголовка `Authorization`
- проверяет токен
- извлекает claims
- сохраняет claims в `context.Context`

### AdminMiddleware

Отвечает за проверку административных прав:

- получает claims из контекста
- проверяет роль пользователя
- разрешает доступ только пользователям с ролью `admin`

### LoggingMiddleware

Логирует информацию о каждом HTTP-запросе:

- HTTP method
- URL path
- HTTP status code
- продолжительность обработки запроса

Пример:

```text
level=INFO msg="http request" method=GET path=/wallets status=200 duration=1.5ms
```

### RecoveryMiddleware

Перехватывает `panic`, возникший во время обработки HTTP-запроса:

- перехватывает panic через `recover`
- записывает ошибку в лог
- возвращает клиенту `500 Internal Server Error`

## context.Context

Контекст HTTP-запроса передаётся через слои приложения:

```text
HTTP Request
    ↓
Handler
    ↓
Service
    ↓
Repository
    ↓
Database
```

Handler передаёт:

```go
r.Context()
```

в service, а service передаёт тот же context в repository.

Repository использует context при выполнении SQL-запросов.

Это связывает операции с базой данных с жизненным циклом конкретного HTTP-запроса.

## Логирование

Для структурированного логирования используется стандартный пакет Go:

```text
log/slog
```

Логируются:

- запуск сервера
- ошибки загрузки конфигурации
- ошибки подключения к базе данных
- HTTP-запросы
- HTTP status code
- время выполнения запроса
- восстановленные panic

## Тестирование

В проекте используются unit-тесты сервисного слоя.

Repository в тестах заменяется fake-реализацией.

Это позволяет тестировать бизнес-логику отдельно от реальной PostgreSQL.

Запуск всех тестов:

```bash
go test ./...
```

Запуск без использования кешированных результатов:

```bash
go test -count=1 ./...
```

## Docker

Проект содержит `Dockerfile` для Go API и `docker-compose.yml` для совместного запуска API и PostgreSQL.

Docker Compose поднимает два сервиса:

```text
api
postgres
```

Внутри Docker-сети API подключается к PostgreSQL по имени сервиса:

```text
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
```

Снаружи PostgreSQL доступен через порт, указанный в `.env`, например:

```text
5433:5432
```

API доступен через:

```text
8080:8080
```

Go-приложение собирается с помощью multi-stage Docker build: на первом этапе создаётся бинарный файл приложения, а в финальный образ попадает только готовый бинарник и минимальное окружение для его запуска.

## Запуск проекта через Docker Compose

### 1. Клонировать репозиторий

```bash
git clone <repository-url>
cd wallets-api-postgres
```

### 2. Создать `.env`

Используйте `.env.example` как шаблон.

Пример:

```env
SERVER_PORT=8080

POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=wallets_db
POSTGRES_SSLMODE=disable

JWT_SECRET=your_secret
```

Файл `.env` содержит секретные данные и не должен попадать в Git или Docker image.

### 3. Собрать и запустить контейнеры

```bash
docker compose up --build
```

После запуска будут работать:

```text
wallets-api
wallets-postgres
```

### 4. Применить миграции

После запуска PostgreSQL примените миграции:

```bash
migrate -path ./migrations -database "postgres://USER:PASSWORD@localhost:5433/wallets_db?sslmode=disable" up
```

### 5. Проверить API

Health check:

```text
http://localhost:8080/health
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

### Остановка проекта

```bash
docker compose down
```

Для удаления контейнеров вместе с volume PostgreSQL:

```bash
docker compose down -v
```

> `docker compose down -v` удалит данные PostgreSQL, хранящиеся в Docker volume.

## Локальный запуск API

PostgreSQL можно оставить запущенным через Docker, а Go API запустить локально:

```bash
go run ./cmd/api
```

В этом случае приложение использует настройки из `.env`:

```env
POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_SSLMODE=disable
```

При запуске API внутри Docker Compose значения подключения к базе переопределяются:

```text
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
```

Таким образом один и тот же Go-код может запускаться как локально, так и внутри Docker.

## Структура проекта

```text
wallets-api-postgres/
├── cmd/
│   └── api/
│       └── main.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── auth/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── response/
│   ├── router/
│   └── service/
├── migrations/
├── tests/
├── .dockerignore
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
└── README_RU.md
```

## Основные изученные темы

В процессе разработки проекта были практически применены:

- REST API на Go
- `net/http`
- PostgreSQL и SQL
- многослойная архитектура
- Repository pattern
- Service layer
- dependency injection через конструкторы
- JWT-аутентификация
- role-based access control
- bcrypt
- middleware
- `context.Context`
- SQL-транзакции
- работа с денежными значениями через `decimal.Decimal`
- SQL-миграции
- структурированное логирование через `slog`
- `defer`, `panic` и `recover`
- unit-тестирование с fake repository
- Swagger / OpenAPI
- Docker
- Docker Compose
- multi-stage Docker build

## Цель проекта

Проект создан для практического изучения backend-разработки на Go и постепенного перехода от простого CRUD API к приложению с разделением ответственности, аутентификацией, авторизацией, бизнес-логикой, транзакциями, тестами, миграциями, документацией и контейнеризацией.