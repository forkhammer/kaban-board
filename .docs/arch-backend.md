# Архитектура backend

## Runtime и wiring

Единственная точка входа — `backend/main.go`. `cmd/api.go` и `cmd/worker.go` являются частями одного процесса, а не отдельными командами.

Composition root находится в `backend/internal/bootstrap/`:

- `di.go` создаёт DB connection, config, memory cache и SprintHub, затем регистрирует repositories, dialect-specific queries, services, use cases и controllers в `goioc/di`;
- injection основан на строковых bean ID и тегах `di.inject`, поэтому переименование ID требует проверки всех consumers;
- `router.go` вручную получает controllers из DI и регистрирует их под `/api`;
- startup вызывает `AutoMigrate` со всеми моделями из `internal/infra/persistance/models/index.go`.

## Слои и направление зависимостей

```text
interfaces/api (Gin controllers, middleware, DTO)
                    |
                    v
app (use cases, services, query/external ports)
                    |
                    v
domain (models, errors, repository ports)
                    ^
                    |
infra (GORM, GitLab, JWT/crypto, cache, WebSocket)
```

Это практическая layered architecture, не строгая Clean Architecture: отдельные application/API типы уже зависят от concrete infra или websocket types. Не расширяйте такие обратные зависимости без необходимости; новые внешние детали закрывайте port-интерфейсом в `app` или `domain`.

## Модули

| Каталог | Ответственность |
| --- | --- |
| `internal/domain/models` | Сущности, enums, validation и поведение accounts, каталога GitLab, sprint/kanban и bindings |
| `internal/domain/repo` | Порты persistence и reports |
| `internal/app/usecases` | CRUD, kanban aggregation, auth/OAuth, reports и GitLab sync orchestration |
| `internal/app/services` | Переиспользуемая application-логика, включая labels и binding history |
| `internal/app/queries` | Filter DTO и query/specification ports |
| `internal/app/interfaces` | Порты GitLab tracker, crypto/password/JWT и SprintHub |
| `internal/interfaces/api/controllers` | Маршруты, bind/validation DTO и преобразование ошибок в HTTP |
| `internal/interfaces/api/middleware` | Optional JWT parsing, route-level auth/admin и online tracking |
| `internal/infra/persistance` | GORM models, repositories и specifications; опечатка `persistance` является реальным именем пакета |
| `internal/infra/db` | Выбор и настройка SQLite/PostgreSQL/MySQL connection |
| `internal/infra/gitlab` | GraphQL pagination и REST GitLab adapter |
| `internal/infra/services` | JWT, password, AES-GCM token encryption и GitLab OAuth |
| `internal/infra/cache`, `internal/infra/hub` | Process-local online cache и WebSocket rooms |

Domain errors `ValidationError`, `NotFoundError` и `ConflictError` отображаются в 400, 404 и 409 через `interfaces/api/utils/errors.go`; неизвестные ошибки отправляются в Sentry и возвращаются как 500.

## HTTP и безопасность

Глобальный `JwtMiddleware` только пытается распознать токен и кладёт account в Gin context. Он не запрещает анонимный запрос. Каждый защищённый маршрут обязан явно подключать `AuthRequired` или `AdminRequired`; UI-проверка роли не является security boundary.

WebSocket endpoint использует SprintHub для binding events. При изменении create/update/delete/order flow сохраняйте согласованность БД, binding history и соответствующего event payload.

## Persistence и запросы

Repositories маппят domain и GORM models. Сложные фильтры строятся через specification implementations в `internal/infra/persistance/spec`. Для issue/project/user/epic существуют отдельные варианты PostgreSQL/MySQL; report repository отдельно специализируется только для PostgreSQL.

При добавлении persistent entity обычно нужны domain model и repository port, ORM model и adapter, регистрация repository в `di.go` и ORM model в `ALL_MODELS`. Добавляйте только действительно задействованные слои, не создавайте пустые abstractions заранее.

Raw SQL обязан учитывать все заявленные dialect. Общий report repository содержит SQLite-style date functions и используется также для MySQL, поэтому report-изменения особенно важно проверять на MySQL.

## GitLab, worker и realtime

- System sync использует private token и GitLab GraphQL/REST adapter.
- OAuth user tokens хранятся в БД зашифрованными AES-GCM; ключ выводится из `ENCRYPTION_KEY`, при пустом значении используется `API_SECRET`.
- Epic представлен GitLab issue с label `GITLAB_EPIC_LABEL`; связи читаются через issue links.
- Worker не допускает overlap собственных ticks, но write-back goroutines не отслеживаются и не ретраятся.
- Cache и SprintHub не durable и не shared между процессами.

## Правила развития

- Новый controller зарегистрируйте и в `registerControllers()` (`di.go`), и в списке `InitRouter()` (`router.go`). Одной DI-регистрации недостаточно.
- Новый repository/query/use case/service регистрируйте с bean ID, точно совпадающим с `di.inject` consumers.
- Новый filter держите в `app/queries`, а GORM implementation в `infra/persistance/spec`; dialect-specific SQL разделяйте рядом с существующими реализациями.
- Многошаговые операции сейчас не получают transaction boundary автоматически. Если нужна атомарность, проектируйте её явно, а не компенсирующими best-effort вызовами.
- Не используйте факт наличия JWT middleware как доказательство авторизации маршрута; проверяйте route group/middleware регистрации controller-а.
- При изменении env обновляйте `config/config.go`, `.env.example` и deployment-конфигурацию вместе.

## Тесты

`backend/tests/integration` поднимает реальный DI/router и `httptest.Server`. SQLite работает in-memory; PostgreSQL/MySQL создаются через Testcontainers и требуют Docker. Worker и реальный GitLab в suite не запускаются.

```bash
# Из корня
mise run test-integration
mise run test-integration-postgresql
mise run test-integration-mysql

# Один test или subtest
cd backend
DB_TYPE=sqlite LOG_LEVEL=1 go test -v -count=1 -timeout 300s ./tests/integration/... -run '^TestGetUsers$'
DB_TYPE=sqlite LOG_LEVEL=1 go test -v -count=1 -timeout 300s ./tests/integration/... -run '^TestGetUsers$/^returns_all_users$'
```

Suite покрывает только часть HTTP API; GitLab sync/client, reports, WebSocket, cache и worker требуют дополнительных targeted tests при изменении.
