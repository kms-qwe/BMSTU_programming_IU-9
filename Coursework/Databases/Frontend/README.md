# Coffee Shop Frontend

Простой учебный frontend для приложения управления кофейней.

Frontend сделан как небольшой SPA без внешних зависимостей и работает поверх backend API по адресу `http://localhost:8080`.

## Что умеет приложение

- проверять наличие активной смены;
- открывать смену;
- закрывать смену;
- показывать список заказов;
- создавать заказы;
- показывать детали заказа;
- показывать и создавать ингредиенты;
- менять активность ингредиентов;
- показывать и создавать пункты меню;
- менять цену пункта меню;
- менять активность пункта меню;
- показывать и создавать складские операции;
- скачивать PDF-отчеты.

## Ограничения и принципы

Frontend реализован строго под переданное OpenAPI backend-а.

Что важно:

- используются только ручки, перечисленные в задании;
- frontend не запрашивает типы складских операций с backend-а;
- ручные типы складских операций захардкожены на фронте:
  - `MANUAL_RESTOCK` — Пополнение
  - `MANUAL_WRITE_OFF` — Списание
- для отображения в списках и отчетах также поддерживаются:
  - `AUTO_ORDER_WRITE_OFF`
  - `INITIAL_STOCK`
- заказы можно только создавать и просматривать;
- ингредиенты и пункты меню не удаляются, а только активируются и деактивируются;
- у пункта меню редактируется только цена;
- название ингредиента, `unit`, название пункта меню и рецепт после создания не редактируются.

## Стек

- `HTML`
- `CSS`
- `JavaScript (ES modules)`
- небольшой `Node.js` static server для SPA-маршрутизации

## Требования для запуска

Нужно, чтобы были доступны:

- `Node.js` версии 18+;
- backend приложения на `http://localhost:8080`.

Отдельный пакетный менеджер не нужен: `npm`, `yarn`, `pnpm` для запуска этого фронта не требуются.

## Запуск

Перейти в директорию frontend:

```bash
cd /Users/mihailkrasnobaev/Desktop/BMSTU_programming_IU-9/Coursework/Databases/Frontend
```

Запустить сервер:

```bash
node server.js
```

После запуска frontend будет доступен по адресу:

- [http://127.0.0.1:3000](http://127.0.0.1:3000)

## Как работает запуск

Файл [server.js](/Users/mihailkrasnobaev/Desktop/BMSTU_programming_IU-9/Coursework/Databases/Frontend/server.js) делает две вещи:

- раздает статические файлы frontend-а;
- возвращает `index.html` для SPA-маршрутов вроде `/orders`, `/menu`, `/inventory`.

Это нужно, чтобы приложение открывалось не только с `/`, но и по прямым ссылкам на страницы.

## Структура файлов

- [index.html](/Users/mihailkrasnobaev/Desktop/BMSTU_programming_IU-9/Coursework/Databases/Frontend/index.html) — корневой HTML-файл приложения
- [styles.css](/Users/mihailkrasnobaev/Desktop/BMSTU_programming_IU-9/Coursework/Databases/Frontend/styles.css) — все стили интерфейса
- [app.js](/Users/mihailkrasnobaev/Desktop/BMSTU_programming_IU-9/Coursework/Databases/Frontend/app.js) — роутинг, API-клиент, состояние приложения, рендер страниц
- [server.js](/Users/mihailkrasnobaev/Desktop/BMSTU_programming_IU-9/Coursework/Databases/Frontend/server.js) — локальный Node.js сервер для раздачи SPA

## Организация страниц

В приложении есть 6 основных frontend-маршрутов.

### `/`

Экран открытия смены.

Логика:

- при загрузке вызывается `GET /api/shift/active`;
- если `is_active = true`, пользователь перенаправляется на `/orders`;
- если `is_active = false`, показывается форма открытия смены;
- форма отправляет `POST /api/shift/open`.

Поля формы:

- `login`
- `password`

### `/orders`

Страница заказов.

Что есть на странице:

- фильтры:
  - сотрудник
  - `from`
  - `to`
- таблица заказов;
- кнопка создания заказа;
- модалка просмотра деталей заказа.

Используемые API:

- `GET /api/order`
- `GET /api/order/{order_id}`
- `POST /api/order`
- `GET /api/menu-item`
- `GET /api/employee`

Особенности:

- при создании заказа frontend отправляет только `menu_item_id` и `quantity`;
- неактивные пункты меню видны в списке, но недоступны для выбора;
- если backend возвращает `pagination`, она отображается на странице.

### `/menu`

Страница пунктов меню.

Что есть на странице:

- список карточек пунктов меню;
- отображение рецепта;
- изменение цены;
- изменение активности;
- создание нового пункта меню.

Используемые API:

- `GET /api/menu-item`
- `GET /api/menu-item/{menu_item_id}`
- `POST /api/menu-item`
- `PATCH /api/menu-item/{menu_item_id}/price`
- `PATCH /api/menu-item/{menu_item_id}/activity`
- `GET /api/ingredient`

Особенности:

- для загрузки списка используется `GET /api/menu-item?include_inactive=true`;
- в рецепте можно выбирать только активные ингредиенты;
- неактивные ингредиенты показываются, но disabled;
- рецепт и название после создания не редактируются.

### `/ingredients`

Страница ингредиентов.

Что есть на странице:

- таблица ингредиентов;
- текущий остаток;
- активность;
- количество рецептов, где используется ингредиент;
- модалка подробной информации;
- создание ингредиента;
- изменение активности.

Используемые API:

- `GET /api/ingredient`
- `GET /api/ingredient/{ingredient_id}`
- `POST /api/ingredient`
- `PATCH /api/ingredient/{ingredient_id}/activity`

Особенности:

- список загружается через `GET /api/ingredient?include_inactive=true`;
- при ошибке `INGREDIENT_USED_IN_RECIPE` показывается сообщение о невозможности деактивации;
- имя и единица измерения после создания не редактируются.

### `/inventory`

Страница складских операций.

Что есть на странице:

- фильтры:
  - сотрудник
  - ингредиент
  - тип операции
  - `from`
  - `to`
- таблица операций;
- создание ручной операции.

Используемые API:

- `GET /api/inventory/operation`
- `POST /api/inventory/operation`
- `GET /api/ingredient`
- `GET /api/employee`

Особенности:

- endpoint получения типов операций не используется;
- типы операций на фронте захардкожены;
- в форме создания доступны только:
  - `MANUAL_RESTOCK`
  - `MANUAL_WRITE_OFF`
- frontend всегда отправляет положительное `amount`;
- знак изменения backend определяет сам.

### `/reports`

Страница PDF-отчетов.

Доступные отчеты:

- отчет по продажам;
- отчет по сотрудникам;
- отчет по складу;
- отчет по популярности меню.

Используемые API:

- `POST /api/report/sales/pdf`
- `POST /api/report/employees/pdf`
- `POST /api/report/inventory/pdf`
- `POST /api/report/menu-popularity/pdf`
- `GET /api/employee`
- `GET /api/ingredient`

Особенности:

- ответы отчетов обрабатываются как binary `blob`;
- frontend не пытается парсить PDF-ответ как JSON;
- файлы скачиваются с именами:
  - `sales-report.pdf`
  - `employees-report.pdf`
  - `inventory-report.pdf`
  - `menu-popularity-report.pdf`

## Layout основного приложения

На защищенных страницах показывается общий layout:

- навигация по разделам;
- информация об активной смене;
- сотрудник, открывший смену;
- время открытия смены;
- кнопка закрытия смены.

Маршруты основного приложения:

- `/orders`
- `/menu`
- `/ingredients`
- `/inventory`
- `/reports`

## Проверка активной смены

Frontend поддерживает следующую логику:

- при открытии `/` выполняется `GET /api/shift/active`;
- при входе на защищенные страницы также выполняется `GET /api/shift/active`;
- если `is_active = false`, пользователь перенаправляется на `/`;
- если любой API-запрос вернул `SHIFT_NOT_ACTIVE`, состояние активной смены очищается и происходит редирект на `/`.

Это соответствует бизнес-правилу, что в системе может быть только одна активная смена, а все пользователи работают внутри нее, пока backend не закроет смену.

## Обработка ошибок

В [app.js](/Users/mihailkrasnobaev/Desktop/BMSTU_programming_IU-9/Coursework/Databases/Frontend/app.js) реализован общий API-клиент.

Он:

- работает с JSON-ответами;
- умеет скачивать binary-ответы для PDF;
- читает `error.code`, `error.message`, `error.details`;
- обрабатывает `SHIFT_NOT_ACTIVE` глобально;
- показывает toast-сообщения для ошибок backend-а.

Поддержаны коды ошибок:

- `SHIFT_NOT_ACTIVE`
- `SHIFT_ALREADY_OPENED`
- `INVALID_CREDENTIALS`
- `DUPLICATE_NAME`
- `INGREDIENT_USED_IN_RECIPE`
- `INACTIVE_INGREDIENT`
- `INACTIVE_MENU_ITEM`
- `INSUFFICIENT_STOCK`
- `UNSUPPORTED_OPERATION`
- `VALIDATION_ERROR`
- `NOT_FOUND`
- `INTERNAL_ERROR`

## UX и поведение интерфейса

В интерфейсе предусмотрены:

- простая навигация;
- таблицы и карточки;
- модальные формы;
- `loading state`;
- `empty state`;
- `error state`;
- `disabled state` для неактивных сущностей;
- бейджи:
  - `Активен`
  - `Неактивен`

## Какие API используются

Frontend использует только эти пути:

### Смены

- `GET /api/shift/active`
- `POST /api/shift/open`
- `POST /api/shift/close`

### Сотрудники

- `GET /api/employee`

### Ингредиенты

- `GET /api/ingredient`
- `GET /api/ingredient/{ingredient_id}`
- `POST /api/ingredient`
- `PATCH /api/ingredient/{ingredient_id}/activity`

### Меню

- `GET /api/menu-item`
- `GET /api/menu-item/{menu_item_id}`
- `POST /api/menu-item`
- `PATCH /api/menu-item/{menu_item_id}/price`
- `PATCH /api/menu-item/{menu_item_id}/activity`

### Заказы

- `GET /api/order`
- `GET /api/order/{order_id}`
- `POST /api/order`

### Складские операции

- `GET /api/inventory/operation`
- `POST /api/inventory/operation`

### Отчеты

- `POST /api/report/sales/pdf`
- `POST /api/report/employees/pdf`
- `POST /api/report/inventory/pdf`
- `POST /api/report/menu-popularity/pdf`

## Какие API не используются

Frontend специально не использует:

- `GET /api/inventory-operation-types`
- любые другие ручки, которых нет в задании

## Что проверить после запуска

Рекомендуемый ручной smoke-test:

1. Открыть [http://127.0.0.1:3000](http://127.0.0.1:3000).
2. Проверить экран открытия смены при отсутствии активной смены.
3. Открыть смену через валидные `login` и `password`.
4. Проверить переход на `/orders`.
5. Создать заказ.
6. Создать ингредиент.
7. Создать пункт меню.
8. Создать ручную складскую операцию.
9. Скачать любой PDF-отчет.
10. Закрыть смену и убедиться, что приложение возвращает на `/`.

## Примечание

Если backend недоступен или не запущен на `http://localhost:8080`, frontend откроется, но запросы будут завершаться ошибками сети.
