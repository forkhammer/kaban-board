# Архитектура системы

## Назначение и состав

Приложение отображает Kanban-доску поверх данных self-hosted GitLab и хранит локальные настройки планирования. Это один репозиторий с двумя независимо собираемыми приложениями:

- `backend/`: Go/Gin REST API, WebSocket endpoint и периодический GitLab sync worker в одном процессе;
- `frontend/board/`: Angular 19 SPA, собираемая в статические файлы и раздаваемая nginx;
- хранилище: SQLite, PostgreSQL или MySQL через GORM;
- внешняя система: GitLab GraphQL для массовой синхронизации и REST v4/OAuth для пользовательских операций.

```text
Browser
  | HTTP /api + WebSocket /api/sprint/:id/ws
  v
Angular SPA ---------------> Gin controllers
                                |
                                v
                         application use cases
                           |             |
                           v             v
                     GORM repositories  GitLab client
                           |             |
                           v             v
                    SQLite/Postgres/MySQL GitLab
                           ^
                           |
                 periodic sync worker (same Go process)
```

Внешнего broker-а или durable queue нет. Online cache и комнаты WebSocket находятся в памяти процесса, поэтому несколько backend replicas не разделяют presence и события.

## Основные потоки

### Запуск backend

`backend/main.go` читает env, инициализирует Sentry и DI, выполняет GORM `AutoMigrate`, запускает sync worker в goroutine, затем блокирующий Gin server. Отдельных API/worker binaries нет.

### Чтение и изменение доски

SPA вызывает REST endpoints с prefix `/api`. Controller преобразует HTTP DTO, use case выполняет бизнес-операцию, repository сохраняет или читает ORM-модель. Изменения issue bindings также формируют history и WebSocket events для открытого спринта.

### Синхронизация GitLab

Worker сразу выполняет sync, затем повторяет его раз в `GITLAB_SYNC_PERIOD_MIN`: projects, releases, users, labels, issues и epics/bindings. Sync последовательный, без общей транзакции, прекращается на первой ошибке и выполняет upsert без удаления исчезнувших GitLab-сущностей.

Некоторые пользовательские изменения сначала сохраняются локально, затем отправляются в GitLab best-effort goroutine. Ошибка внешней записи логируется, но локальную транзакцию не откатывает.

## Данные

ORM-модели перечислены в `backend/internal/infra/persistance/models/index.go`; этот список одновременно задаёт набор `AutoMigrate`. Основные агрегаты: accounts, teams/groups/users/projects, labels, issues/epics/releases, sprints/columns, issue bindings и их history, sprint user settings, key-value settings и encrypted GitLab OAuth tokens.

Отдельных versioned migrations нет. Изменение схемы вступает в силу при следующем старте через `AutoMigrate`; destructive преобразования данных автоматически не выполняются. SQL/search specifications различаются по dialect, поэтому поддержку трёх БД нельзя считать эквивалентной без трёх integration-запусков.

## Развёртывание и конфигурация

- `docker-compose.yml`: backend + frontend + SQLite volume.
- `docker-compose.postgresql.yml` и `docker-compose.mysql.yml`: расширяют базовый compose соответствующей БД.
- Backend получает env из корневого `.env`; полный актуальный набор полей определён в `backend/config/config.go`, а `.env.example` покрывает не все OAuth/Sentry/encryption параметры.
- Frontend development использует `http://localhost:8080/api`. Production читает `window.env.apiUrl` из `assets/env.js`, иначе использует same-origin `/api`.
- Frontend Dockerfile создаёт `env.js` из `API_URL` во время build, а не при запуске контейнера. nginx раздаёт только SPA и не proxy-ирует `/api`.
- Compose передаёт frontend build arg `API_URL=http://localhost:8080/api`; при доступе с другой машины этот localhost относится к браузеру пользователя.

## Операционные ограничения

- Compose `depends_on` не ждёт готовности PostgreSQL/MySQL, а backend не повторяет подключение к БД.
- У backend нет graceful shutdown по OS signals; отмена worker context происходит только после возврата Gin server.
- Process-local cache/WebSocket hub ограничивают горизонтальное масштабирование без отдельного shared transport/state.
- GitLab sync и write-back не образуют распределённую транзакцию; при доработках явно проектируйте поведение при частичном отказе.
