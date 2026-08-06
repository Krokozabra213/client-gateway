# client-gateway

HTTP API-шлюз и балансировщик платформы Gargantua Labs. Единая точка входа для клиентов, маршрутизирующая запросы к внутренним микросервисам через gRPC, обрабатывающая Telegram-вебхуки ботов и обеспечивающая защиту от перегрузок с помощью rate limiter на основе алгоритма Token Bucket.

## Назначение

Сервис предоставляет REST API для клиентских приложений, управляет регистрацией и аутентификацией пользователей через Telegram, обеспечивает подключение и управление пользовательскими Telegram-ботами, обрабатывает пополнения баланса (Telegram Stars, TBank), покупку продуктов и подписок, а также принимает уведомления от платёжной системы TBank.

Кроме того, client-gateway выступает в роли **балансировщика нагрузки**: распределяет входящие вебхуки от множества подключённых пользовательских ботов, маршрутизирует их к соответствующим обработчикам и ограничивает интенсивность запросов через **Rate Limiter** (алгоритм Token Bucket), реализованный на базе Redis — для каждого бота ведётся подсчёт RPS, что предотвращает превышение допустимой нагрузки на внутренние сервисы.

## Технологический стек

- **Язык**: Go 1.25
- **HTTP-фреймворк**: [Fiber v3](https://github.com/gofiber/fiber)
- **Telegram Bot API**: [telebot.v4](https://github.com/mymmrac/telego)
- **База данных**: PostgreSQL 16 (через pgx/v5)
- **Кэш**: Redis 8
- **Межсервисное взаимодействие**: gRPC
- **Наблюдаемость**: OpenTelemetry (traces + metrics), Prometheus
- **Контейнеризация**: Docker, docker-compose с профилями (`prod`, `test`)

## Архитектура

Сервис построен по модульному принципу (feature-based architecture). Каждый модуль независим и содержит собственные слои: handler → service → repository.

```
cmd/client-gateway/main.go        # точка входа, DI-контейнер
configs/                           # конфигурации (dev, stage, prod)
internal/
  ├── config/                      # структуры конфигурации
  ├── clients/                     # gRPC-клиенты (payment, balance, order-product)
  ├── domain/                      # общие доменные модели
  ├── security/                    # сервисы шифрования и хеширования
  ├── server/http/                 # HTTP-сервер (Fiber)
  └── features/
      ├── auth/                    # модуль аутентификации
      │   ├── domain/              # User, UserID, Username, UserRole, UserGender
      │   ├── handler/             # HTTP-обработчики
      │   ├── repository/          # postgres + redis
      │   └── service/             # Register, GetUserByTelegramID
      ├── bot/                     # модуль Telegram-ботов
      │   ├── constants/           # callback-ы, команды, разметки
      │   ├── domain/              # Bot, BotID
      │   ├── error_handler/       # обработчик ошибок бота
      │   ├── handler/             # обработчики кнопок и команд бота
      │   ├── middlewares/         # authenticate, owner_only, rps
      │   ├── repository/          # postgres + redis
      │   └── service/             # ConnectBot, GetConnectedBot, GetBotsByUserID...
      ├── shared/                  # общие типы и ошибки
      │   └── types/               # NonEmptyString, PositiveInt, TelegramID, Pagination
      └── tbank/                   # модуль уведомлений TBank
          └── handler/             # HandleNotification
migrations/                        # SQL-миграции (goose)
```

## API Endpoints

| Метод | Путь | Назначение |
|-------|------|------------|
| `GET` | `/metrics` | Prometheus-метрики |
| `GET` | `/health` | Health-check (Redis + PostgreSQL) |
| `POST` | `/v1/auth/...` | Роуты аутентификации (в разработке) |
| `POST` | `/v1/bot/` | Внутренний вебхук платформенного бота |
| `POST` | `/v1/bot/:botId` | Внешний вебхук подключённого бота пользователя |
| `POST` | `/v1/tbank/notification` | Уведомления о платежах TBank |

## Модули

### Auth

Управление пользователями платформы:

- `POST /v1/auth/...` — регистрация и аутентификация (через Telegram)
- Хранение: PostgreSQL (основное) + Redis (кэш с TTL)
- Поля пользователя: `id`, `telegram_id`, `first_name`, `username`, `birth_date`, `gender`, `role`, `is_verified`

### Bot

Управление Telegram-ботами через платформенного бота:

- Регистрация пользователя (выбор пола, даты рождения)
- Главное меню: баланс, боты, магазин, заказы, подписки
- Подключение бота: ввод токена → валидация → сохранение (токен шифруется)
- Просмотр списка ботов с пагинацией, информация о боте
- Пополнение баланса: Telegram Stars (через `payment.TGInit`) или TBank (через `payment.TbankInit`)
- Обработка Telegram Checkout (`OnCheckout`, `OnPayment`)
- Каталог товаров (`order-product.ListProducts`), просмотр продукта, покупка
- Заказы (`order-product.ListOrders`) и подписки (`order-product.ListSubscriptions`)

Мiddleware бота:
- `Authenticate` — проверяет авторизацию пользователя по `telegram_id`
- `OwnerOnly` — ограничивает доступ владельцем бота
- `CalculateBotRPS` — подсчитывает RPS бота в Redis

### Tbank

Приём уведомлений от TBank о статусе платежей:

- `POST /v1/tbank/notification` — webhook, пробрасывает запрос в `payment.TbankWebhook`

## База данных

| Таблица | Назначение |
|---------|------------|
| `users` | Пользователи платформы |
| `bots` | Подключённые Telegram-боты (токен зашифрован, token_hash для уникальности) |

Миграции управляются через [goose](https://github.com/pressly/goose).

## gRPC-зависимости

| Сервис | Назначение |
|--------|------------|
| **payment** (`:44050`) | Инициация платежей (TG Stars, TBank), webhook TBank |
| **balance** (`:44070`) | Получение баланса аккаунта |
| **order-product** (`:44030`) | Каталог продуктов, заказы, подписки |

## Переменные окружения

Конфиденциальные переменные окружения хранятся в зашифрованном виде в файле `.env.vault` (Ansible Vault) и не попадают в репозиторий в открытом виде.

### Получение `.env` из `.env.vault`

Для расшифровки `.env.vault` в `.env` используется Ansible Vault:

```bash
# Расшифровка .env файла
ansible-vault view .env.vault > .env
# Пароль: 7d28a604-bb66-4ba9-a1ed-7679a2f82001

# Просмотр без сохранения
ansible-vault view .env.vault

# Обратное шифрование после изменений
ansible-vault encrypt .env --output .env.vault
```

После расшифровки файл `.env` загружается автоматически через `godotenv.Load()` при старте сервиса.

Основные переменные в `.env`:

| Переменная | Назначение |
|------------|------------|
| `ENV` | Окружение (`dev`, `stage`, `prod`), по умолчанию `dev` |
| `USER_CACHE_TTL_MINUTES` | TTL кэша пользователя |
| `ENCRYPTION_KEY` | Ключ шифрования |
| `HASHING_KEY` | Ключ хеширования |
| `BOT_TOKEN` | Токен платформенного бота |
| `BOT_WEBHOOK_URL` | Базовый URL для вебхуков ботов |
| `BOT_SECRET` | Секрет для вебхуков |

## Запуск

### Локально

```bash
cd client-gateway
cp .env.example .env   # заполнить переменные
make                    # или go run ./cmd/client-gateway
```

### Docker (prod-профиль)

```bash
make docker-prod
# или
docker compose --profile prod up -d
```

Поднимаются: `api`, `postgres`, `redis`, `migrate`, `otel-collector`.

### Тесты

```bash
# unit-тесты
make test

# интеграционные тесты (требуют test-postgres)
make docker-test-up
make test-integration
make docker-test-down
```

## CI/CD

- **CI** (`.github/workflows/ci.yaml`): линтинг и тесты при push/PR
- **Deploy** (`.github/workflows/deploy.yaml`): деплой при push в ветки `main`, `develop`, `stage`

## Health-check

`GET /health` возвращает:

```json
{
  "status": "ok",
  // или при ошибке:
  "status": "error",
  "redis": "down",    // если Redis недоступен
  "postgres": "down"   // если PostgreSQL недоступен
}
```

HTTP-статус: `200` (OK) или `500` (есть проблемы).
