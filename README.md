# Market

Учебный маркетплейс на Go. В проекте есть отдельные сервисы для пользователей, товаров и корзины, а также простой фронтенд.

## Состав проекта

| Компонент | Назначение | Порт |
| --- | --- | --- |
| `auth-service` | Регистрация, вход и информация о пользователе | `8080` |
| `product-service` | Список товаров и получение товара по ID | `8081` |
| `cart-service` | Корзина авторизованного пользователя | `8082` |
| `frontend` | Каталог, регистрация, вход и корзина | `3000` |
| PostgreSQL | Пользователи и товары | внутри Compose |
| Redis | Корзины и кеш товаров | внутри Compose |

Go-код находится в `backend/`, фронтенд — в `frontend/`. Точки входа сервисов лежат в `backend/cmd/`, логика — в `backend/internal/`.

## Запуск через Docker Compose

Из корня проекта:

```bash
docker compose up --build -d
```

Откройте [http://localhost:3000](http://localhost:3000). Для просмотра состояния и остановки:

```bash
docker compose ps
docker compose down
```

При первом запуске PostgreSQL создаёт таблицы и тестовые товары из файлов в `backend/internal/*/migrations/`. Данные PostgreSQL и Redis сохраняются в Docker volumes. Для Redis включён AOF. Миграции при следующих запусках не выполняются повторно, пока существует том PostgreSQL.

Оба сервиса, работающие с токенами, должны использовать одинаковый `TOKEN_SECRET`. В Compose можно передать его через переменную окружения:

```powershell
$env:TOKEN_SECRET = "your-secret"
docker compose up --build -d
```

Если переменная не задана, используется учебное значение `secret`.

## API

### Авторизация — `localhost:8080`

| Метод и путь | Назначение |
| --- | --- |
| `POST /auth/register` | Создать пользователя |
| `POST /auth/login` | Получить access token |
| `GET /auth/me` | Получить текущего пользователя; нужен Bearer-токен |

Пример регистрации:

```json
{"username":"alice","password":"password123","account_type":"user"}
```

После регистрации нужно войти через `POST /auth/login` с `username` и `password`. Токен передаётся в заголовке `Authorization: Bearer <access_token>`. Срок его действия — 3 минуты.

### Товары — `localhost:8081`

| Метод и путь | Назначение |
| --- | --- |
| `POST /products/list` | Получить список товаров |
| `GET /products?id=1` | Получить товар по ID |

Тело запроса для списка:

```json
{"number":20,"p_types":["sport","books"]}
```

Чтобы получить товары всех категорий, передайте пустой массив `p_types`. Доступные категории в тестовых данных: `sport`, `electronics`, `books`, `clothes`.

### Корзина — `localhost:8082`

Все запросы к корзине требуют Bearer-токен.

| Метод и путь | Назначение | Тело запроса |
| --- | --- | --- |
| `GET /cart` | Получить корзину | — |
| `POST /cart/items` | Добавить количество к товару | `{"product_id":1,"quantity":2}` |
| `PUT /cart/items/1` | Установить количество товара | `{"quantity":3}` |
| `DELETE /cart/items/1` | Удалить товар | — |
| `DELETE /cart` | Очистить корзину | — |

Корзина хранит ID товара и количество. Название и цену фронтенд получает из `product-service`; проверка наличия товара и оформление заказа пока не реализованы.

## Локальный запуск без Compose

Нужны Go версии из `backend/go.mod`, Node.js, PostgreSQL и Redis. Создайте базу `market` и примените SQL-файлы из `backend/internal/auth/migrations/` и `backend/internal/product/migrations/`. Затем запустите сервисы в отдельных терминалах:

```bash
cd backend
go run ./cmd/auth-service
go run ./cmd/product-service
go run ./cmd/cart-service
```

Для фронтенда из корня проекта:

```bash
node frontend/server.js
```

Адреса подключений настраиваются через `DATABASE_URL`, `REDIS_ADDR`, `TOKEN_SECRET`, `PRODUCT_SERVICE_URL`, `AUTH_SERVICE_URL` и `CART_SERVICE_URL`. Фронтенд передаёт запросы `/products`, `/auth` и `/cart` соответствующим сервисам, поэтому браузеру не нужна отдельная настройка CORS.
