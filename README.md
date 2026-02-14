# market-parser

Парсер товаров и цен на Go.

Проект собирает список товаров по категориям из онлайн-магазина и возвращает JSON с полями `name`, `price`, `link`.

> В проекте реализован парсер для выбранного магазина с запуском из `Makefile` / `docker-compose`.

---

## Что реализовано

* Парсер на Go (без GUI, headless Chromium через rod).
* Поддержка выбора магазина, адреса доставки и категории товаров.
* Возможность использовать прокси для браузера (настройки через `.env`).
* Серверный режим: HTTP API `/api/v1/market-parser/parse`.
* Конфигурация через `configs/config.yaml` и переменные окружения (.env.example).

## Быстрый обзор структуры (файлы/директории важные для запуска)

* `cmd/market-parser/main.go` — точка входа, инициализация конфигурации и сервисов.
* `configs/config.yaml` — основной YAML-конфиг (base_url, selectors и т.д.).
* `.env.example` — пример переменных окружения (WebSocket URL браузера, прокси и т.п.).
* `internal/adapters/browser/chromium` — инициализация Chromium/rod, прокси, DTO и утилиты для работы с элементами страницы.
* `internal/adapters/parsers/kuper.go` — пример реализации парсера для магазина.
* `internal/usecase/parser_service.go` — сервисная логика: валидация параметров и вызов репозитория-парсера.
* `internal/transport/http` — HTTP-обёртка, handlers и OpenAPI-строка.
* `Dockerfile`, `docker-compose.yaml`, `Makefile` — способы запуска.

---

## Требования

* Go 1.20+ (или версия, указанная в `go.mod`).
* Docker & docker-compose (для контейнерного запуска).
* Локальный Chromium (для локального запуска через Makefile) или контейнер `chromium` (для docker-compose).

---

## Конфигурация (.env)

Скопируйте `.env.example` в `.env` и заполните значения:

* `BROWSER_WS_URL` — WebSocket адрес браузера (например `ws://localhost:7317` или `ws://chromium:7317` в docker-compose).
* `BROWSER_HEADLESS` — `true`/`false`.
* `BROWSER_PROXY_HAS` — `true`/`false` (включить использование прокси).
* `BROWSER_PROXY_IP`, `BROWSER_PROXY_PORT`, `BROWSER_PROXY_LOGIN`, `BROWSER_PROXY_PASSWORD` — параметры прокси.

Также, в `config.yaml` имеются поля:

* `headless_mode` — `true`/`false` (headful/headless).
* `test_mode` — `true`/`false` (для быстрого тестирования функционала парсинга страниц).
* `human_like_mode` — `true`/`false` (вкл./вык. автоматическое движение мыши/скроллинг).

---

## Способ A — локальный запуск (с локальным Chromium через Makefile)

Этот способ удобен для разработки и отладки (вы видите браузер). В Makefile есть цель для запуска приложения и локального браузера.


1. Подготовьте `.env` на основе `.env.example` и `config.yaml`.

2. Соберите и запустите приложение через Makefile:

```bash
# сборка
make build

# запуск (запускает приложение и подключается к локальному Chromium)
make run

# запустить swagger ui по адресу http://localhost:8081 для отправки запросов
make swaggerui
```

3. Пример HTTP-запроса к API (парсинг категории):

```bash
curl 'http://localhost:8080/api/v1/market-parser/parse?category=Мясо, птица&address=Москва, Красная площадь, 3&market=metro'
```

Ответ — JSON-массив объектов `{ "name": "...", "price": 123.0, "link": "https://..." }`.

---

## Способ B — запуск через Docker Compose

Контейнерный запуск рекомендован для воспроизводимости и тестового деплоя. В `docker-compose.yaml` описаны два сервиса: `market-parser` и `chromium`.

1. Подготовьте `.env` на основе `.env.example` и `config.yaml`.

2. Запустите compose:

```bash
docker-compose up --build
```

3. После поднятия сервисов API будет доступен на порту, указанном в `configs/config.yaml` / `.env` (по умолчанию `localhost:8080`).

> Важно: по умолчанию парсер в контейнере запускается в headless-режиме — убедитесь, что `headless_mode: true`.

---

## Пример использования API

**GET** `/api/v1/market-parser/parse`

Параметры:

* `market` — идентификатор магазина/профиля (обязательный).
* `address` — адрес доставки / профиль (обязательный).
* `category` — категория товаров (обязательный).

Пример:

```bash
curl 'http://localhost:8080/api/v1/market-parser/parse?category=Мясо, птица&address=Москва, Красная площадь, 3&market=metro'
```

* Также можете воспользовать swagger ui (http://localhost:8081 - при запуске через docker-compose).

Ответ:

```json
[
  { "name": "Куриное филе", "price": 371.0, "link": "https://..." },
  { "name": "Свинная шейка", "price": 456.0, "link": "https://..." }
]
```

---

## Настройка прокси

Проект поддерживает прокси: параметры передаются через `.env` и применяются при инициализации Chromium (WebSocket соединение конфигурируется с `Proxy` и при необходимости вызывается `HandleAuth`).

---
