# wallet-service
Сервис, разработанный на Go, предназначенный для управления балансами пользователей (кошельками).
## Технологический стек
Язык: 
- Go (Golang)
- HTTP-фреймворк: go-chi/chi
- База данных: PostgreSQL
- ORM/DB: sqlx (с драйвером pgx)
- Точные числа: decimal
- Развертывание: Docker, Docker Compose
## Структура проекта
```
wallet-service
├── Dockerfile                      # сборка докер образа для сервиса
├── README.md
├── cmd
     └── wallet
         └── main.go                # корень приложения
├── config.env                      # переменные среды для сервера, бд, и миграции
├── docker-compose.yml              # поднимает весь проект с бд
├── go.mod
├── go.sum
└── internal                        # папка с внутренними для проекта пакетами
├── app
     └── app.go                     # первая точка входа в приложение
├── config
     └── config.go                  # структура, загружающая конфигурацию из config.env
├── handlers
     └── handlers.go                # обработчики для endpoint'ов сервиса
     └── handlers_test.go
├── migrations
     └── 001_create_wallet.up.sql   # стартовая миграция
├── repo
     ├── mock_repo.go       
     └── repo.go                    # логика взаимодействия с бд                    
├── service
     ├── mock_service.go
     ├── service.go                 # бизнес-логика
     └── service_test.go
└── wallet
    └── wallet.go                   # структура wallet
```
## Запуск проекта
### Требования
Для запуска проекта необходимы:
- Установленный Docker и Docker Compose.
- Go 1.25.5+.
### Сборка и запуск
Выполните следующие команды в корневой директории проекта:
Сборка образов и запуск контейнеров:
```bash
  docker-compose up --build -d
```
Для удаления контейнера и связанных с ним томов используйте:
```bash
  docker-compose down -v
```
### Переменные окружения
Конфигурация сервиса задается в файле config.env:
- DB_HOST - Имя хоста БД (в Docker-сети)
- DB_PORT - Порт БД (в Docker-сети)
- DB_PASSWORD - Пароль пользователя 
- MIGRATE_PATH - Абсолютный путь к миграциям внутри контейнера
- SERVER_ADDR - Адрес прослушивания внутри контейнера
- SERVER_PORT - Порт прослушивания
## Использование API 
API использует стандартные HTTP-методы и кодировку JSON. Порт для доступа с хоста — 9000.
1. Пополнение / Снятие средств (Change Amount)
Создает кошелек, если он не существует (при первом DEPOSIT большем 0), или изменяет существующий баланс.
Пример запроса:
```bash
  curl -X POST http://localhost:9000/api/v1/wallet -H "Content-Type: application/json" -d '{ "wallet_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef", "operation_type": "DEPOSIT", "amount": 1000 }'
```
Пример снятия средств:
```bash
  curl -X POST http://localhost:9000/api/v1/wallet -H "Content-Type: application/json" -d '{ "wallet_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef", "operation_type": "WITHDRAW", "amount": 1000 }'
```
2. Получение баланса (Get Amount)
Возвращает текущий баланс для указанного кошелька
```bash
  curl -X GET http://localhost:9000/api/v1/wallets/a1b2c3d4-e5f6-7890-1234-567890abcdef
```
