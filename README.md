# 🚀 Быстрый запуск с Docker Compose

## Шаги

```bash
# Клонирование репозитория
git clone ...
cd ...

# Запуск контейнеров
docker compose up --build
```

# 📄 Описание переменных окружения

# Переменные для manager-сервиса
```bash
MANAGER_SERVER_HOST           # Адрес, на котором работает manager-сервис
MANAGER_SERVER_PORT           # Порт HTTP-сервера manager-сервиса
MANAGER_METRICS_PORT          # Порт для экспорта метрик manager-сервиса
MANAGER_RETRY_INTERVAL_SEC    # Интервал между попытками переподключения к node-сервисам (в секундах)
MANAGER_RECOVERY_INTERVAL_SEC # Интервал между попытками восстановления node-сервисов (в секундах)
```

```bash
# Переменные для node-сервисов
NODE_SERVER_HOST              # Адрес, на котором работает node-сервис
NODE_SERVER_PORT              # Порт HTTP-сервера node-сервиса
NODE_METRICS_PORT             # Порт для экспорта метрик node-сервиса
NODE_NAME                     # Уникальное имя node-сервиса
```

```bash
# Общие переменные
GRPC_MANAGER_HOST             # Адрес manager-сервиса для gRPC-соединения
```

```bash
# Переменные Grafana
GRAFANA_ADMIN_PASSWORD        # Пароль администратора Grafana
GRAFANA_PORT                  # Порт для доступа к интерфейсу Grafana
```

```bash
# Переменные Prometheus
PROMETHEUS_MAX_DAYS           # Максимальное количество дней для хранения данных Prometheus
PROMETHEUS_MAX_SIZE_GB        # Максимальный объем данных Prometheus (в ГБ)
PROMETHEUS_PORT               # Порт для доступа к интерфейсу Prometheus
```
