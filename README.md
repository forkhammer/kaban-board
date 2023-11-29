# Kanban board for self-hosted gitlab

Возможности доски:
 - Быстрая загрузка пользователей и тикетов
 - Фоновая сихронизация с Gitlab
 - Поиск по тикетам и пользователям
 - Тонкая настройка колонок
 - Группы проектов
 - Темная/светлая тема
 - Встроенная административная панель


## Start

### Docker

### Local environment

### Configuration

Environment variables

HOST=0.0.0.0

PORT=8080

GITLAB_URL=

GITLAB_TOKEN=


POSTGRES_HOST=localhost

POSTGRES_PORT=5432

POSTGRES_DB=board

POSTGRES_USER=board

POSTGRES_PASSWORD=board

LOG_LEVEL=1

JWT_TOKEN_LIFESPAN_HOUR=24

API_SECRET=JWT_SECRET

GITLAB_SYNC_PERIOD_MIN=10

MEMORY_CACHE_DURATION_MIN=15