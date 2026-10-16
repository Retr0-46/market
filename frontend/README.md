# Frontend

Простая витрина с регистрацией, входом и корзиной без дополнительных пакетов.

1. Запустите PostgreSQL, Redis, auth-service, product-service и cart-service. Проще всего сделать это через `docker compose up --build -d` из корня проекта.
2. В отдельном терминале запустите `node frontend/server.js` из корня проекта.
3. Откройте http://localhost:3000.

Сервер фронтенда передаёт запросы `/products`, `/auth` и `/cart` соответствующим сервисам. Локальные адреса по умолчанию: `localhost:8081`, `localhost:8080` и `localhost:8082`. Их можно изменить переменными `PRODUCT_SERVICE_URL`, `AUTH_SERVICE_URL`, `CART_SERVICE_URL`; порт страницы — через `FRONTEND_PORT`.

Токен входа хранится только в памяти вкладки. После обновления страницы потребуется войти снова. Корзина остаётся в Redis.
