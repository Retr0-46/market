# Frontend

Простая витрина для product-service без дополнительных пакетов.

1. Запустите PostgreSQL, Redis и product-service (`cd backend; go run ./cmd/product-service`).
2. В отдельном терминале запустите `node frontend/server.js` из корня проекта.
3. Откройте http://localhost:3000.

Сервер фронтенда передаёт запросы `/products` в product-service на `http://localhost:8081`. Адреса можно изменить переменными `FRONTEND_PORT` и `PRODUCT_SERVICE_URL`.
