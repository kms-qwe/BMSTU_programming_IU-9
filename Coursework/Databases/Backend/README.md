# Coffee Shop Backend

Go-бэкенд кофейни с локальным запуском через `docker compose`.

## Структура

```text
.
├── app
│   ├── cmd
│   ├── internal
│   └── pkg
├── build
├── docs
├── migrations
├── scripts
├── Makefile
└── README.md
```

## Как запустить

1. Подготовить env:

```bash
cp build/.env.example build/.env
```

2. Поднять PostgreSQL, миграции и приложение:

```bash
make up
```

3. Дождаться старта сервиса.

По умолчанию API будет доступен по адресу:

```text
http://localhost:8080
```

Остановка:

```bash
make down
```

## Команды Makefile

Все основные команды:

- `make run`:
  запускает приложение локально из `app/cmd/app`
- `make fmt`:
  форматирует Go-код в `app/cmd`, `app/internal`, `app/pkg`
- `make tidy`:
  запускает `go mod tidy` в папке `app`
- `make install-goose`:
  скачивает `goose` в локальную директорию `.localbin`
- `make install-swag`:
  скачивает `swag` в локальную директорию `.localbin`
- `make install-tools`:
  скачивает сразу `goose` и `swag`
- `make swagger`:
  генерирует swagger-документацию из комментариев в `docs/swagger`
- `make migration-create name=<migration_name>`:
  создает новую SQL-миграцию через `goose` в папке `migrations`
- `make up`:
  поднимает PostgreSQL, миграции и приложение через `docker compose`
- `make down`:
  останавливает контейнеры и удаляет volume базы данных

Примеры:

```bash
make install-tools
make migration-create name=add_orders_indexes
make swagger
make up
```

## Локальный запуск без Docker

Если PostgreSQL уже поднят отдельно и переменные окружения выставлены:

```bash
cd app
go run ./cmd/app
```

## Как создать сотрудника

Для быстрого добавления сотрудника в уже поднятую БД есть скрипт:

```bash
./scripts/create_employee.sh "Иван Петров" "+79990000000" "ivan" "123456"
```

Скрипт использует настройки из `build/.env` и выполняет SQL внутри контейнера с PostgreSQL.

## Как слать запросы

Примеры ниже предполагают, что сервис поднят на `http://localhost:8080`.

Проверить активную смену:

```bash
curl http://localhost:8080/api/shift/active
```

Открыть смену:

```bash
curl -X POST http://localhost:8080/api/shift/open \
  -H "Content-Type: application/json" \
  -d '{
    "login": "ivan",
    "password": "123456"
  }'
```

Получить сотрудников:

```bash
curl http://localhost:8080/api/employee
```

Создать ингредиент:

```bash
curl -X POST http://localhost:8080/api/ingredient \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Молоко",
    "unit": "л",
    "initial_stock": 10
  }'
```

Создать заказ:

```bash
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "menu_item_id": "54c2ec90-9a67-41fd-84dc-891d527e5482",
        "quantity": 2
      }
    ]
  }'
```

Скачать Excel-отчет:

```bash
curl -X POST http://localhost:8080/api/report/sales/excel \
  -H "Content-Type: application/json" \
  -d '{
    "from": "2026-04-01T00:00:00Z",
    "to": "2026-04-26T23:59:59Z",
    "employee_id": null
  }' \
  --output sales-report.xlsx
```

## Документация API

Полная OpenAPI-документация лежит в:

- [docs/openapi.yaml](/Users/mihailkrasnobaev/Desktop/BMSTU_programming_IU-9/Coursework/Databases/Backend/docs/openapi.yaml)

Список доменов:

- `shift`
- `employee`
- `ingredient`
- `menu-item`
- `order`
- `inventory`
- `report`

## Swagger из комментариев

В handler’ы добавлены swagger-комментарии. Для генерации swagger-файлов можно использовать:

```bash
make swagger
```

Команда генерирует документацию в `docs/swagger`.

## Миграции

Создать новую миграцию:

```bash
make migration-create name=create_suppliers_table
```

После этого в папке `migrations` появятся новые `.up.sql` и `.down.sql` файлы, которые можно заполнить SQL-кодом.

## Все API методы

### Shift

- `GET /api/shift/active`
- `POST /api/shift/open`
- `POST /api/shift/close`

### Employee

- `GET /api/employee`

### Ingredient

- `GET /api/ingredient`
- `GET /api/ingredient/{ingredient_id}`
- `POST /api/ingredient`
- `PATCH /api/ingredient/{ingredient_id}/activity`

### Menu Item

- `GET /api/menu-item`
- `GET /api/menu-item/{menu_item_id}`
- `POST /api/menu-item`
- `PATCH /api/menu-item/{menu_item_id}/price`
- `PATCH /api/menu-item/{menu_item_id}/activity`

### Order

- `GET /api/order`
- `GET /api/order/{order_id}`
- `POST /api/order`

### Inventory

- `GET /api/inventory/operation`
- `POST /api/inventory/operation`

### Report

- `POST /api/report/sales/excel`
- `POST /api/report/employees/excel`
- `POST /api/report/inventory/excel`
- `POST /api/report/menu-popularity/excel`
