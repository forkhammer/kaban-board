# Kanban board for self-hosted gitlab

Возможности доски:
 - Быстрая загрузка пользователей и тикетов
 - Фоновая сихронизация с Gitlab
 - Поиск по тикетам и пользователям
 - Тонкая настройка колонок
 - Группы проектов
 - Темная/светлая тема
 - Встроенная административная панель


 ![Kaban board](/board.png)


## Start

### Docker

Running through Docker-containers

```
docker-compose up -d
```

### Local environment

Running in local environment

```
./scripts/run_dev_backend.sh

./scripts/run_dev_frontend.sh
```


### Configuration

**Environment variables**

You can set local environment settings through the file **.env.** Copy **.env.example** to **.env**

HOST=0.0.0.0   # backend host

PORT=8080   # backend port

GITLAB_URL=     # set your gitlab url. Examle: https://gitlab.yourdomain.com

GITLAB_TOKEN=       # set gitlab private token for api access

DB_TYPE=sqlite      # type of database for storing settings: sqlite, postgresql, mysql


POSTGRES_HOST=localhost

POSTGRES_PORT=5432

POSTGRES_DB=board

POSTGRES_USER=board

POSTGRES_PASSWORD=board


MYSQL_HOST=localhost

MYSQL_PORT=3306

MYSQL_DATABASE=board

MYSQL_ROOT_PASSWORD=board

MYSQL_USER=board

MYSQL_PASSWORD=board


SQLITE_DB_FILE=sqlite.db


LOG_LEVEL=1     # log level: 1 - silent, 2 - error, 3 - warn, 4 - info

JWT_TOKEN_LIFESPAN_HOUR=24      # JWT-token lifetime in hours

API_SECRET=JWT_SECRET   # Secret key for JWT-token generation

GITLAB_SYNC_PERIOD_MIN=10  #  Period for gitlab sync in minutes

MEMORY_CACHE_DURATION_MIN=15    # Local memory cache lifetime

ALLOW_ORIGINS=http://localhost:4200     # Origins for access to the backend. Can be listed separated by commas