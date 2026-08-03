# Правила кода и слоёв

Общие обязательные принципы находятся в корневом `AGENTS.md`. Здесь перечислены только правила, зависящие от устройства этого репозитория.

## Границы изменений

- Сначала проследите существующий end-to-end поток от route/component до use case/service и persistence/API. Меняйте минимальный набор реально участвующих слоёв.
- Не переносите transport DTO в domain и не протаскивайте GORM/Angular implementation details через application contracts.
- Существующее нарушение границ не является шаблоном для нового кода. Новая интеграция должна зависеть от узкого port, если это отделяет business logic от внешней библиотеки.
- Не создавайте пустые repository/service/facade abstractions «на будущее». Слой появляется вместе с реальным поведением и consumer-ом.

## Backend

- Controllers отвечают за HTTP parsing/DTO/status, use cases за orchestration и бизнес-правила, domain models за invariants, repositories/specifications за persistence/query details.
- Возвращайте типизированные domain errors, когда endpoint должен дать 400/404/409; неизвестная ошибка проходит через `utils.HandleException` как 500/Sentry.
- Защищайте mutation routes явным `AuthRequired`/`AdminRequired`. Глобальный JWT middleware только заполняет context.
- При новой ORM-модели обновляйте `models.ALL_MODELS`; иначе `AutoMigrate` её не создаст.
- При новом controller обновляйте одновременно DI registration и ручной список router-а.
- Проверяйте query/raw SQL на SQLite, PostgreSQL и MySQL. Не добавляйте общий SQL с dialect-specific date/search syntax.
- Если операция должна быть атомарной, добавьте реальную transaction boundary. Не маскируйте частичный commit fallback-ом или фоновой компенсацией.
- Изменения issue binding обязаны сохранять history и согласованный Sprint WebSocket event.
- Не исправляйте написание каталога `persistance` частично: это package path по всему backend, а не локальная опечатка.
- Go-код форматируйте `gofmt`; для нового кода используйте `any`, как требует текущая версия Go и существующий стиль.

## Frontend

- Соблюдайте текущую границу: route pages в `app/components`/standalone reports, domain feature в `modules/kanban`, reusable controls в `modules/ui`, singleton infrastructure в `modules/core`.
- Не импортируйте `CoreModule.forRoot()` вне `AppModule`.
- Для NgModule component используйте `standalone: false` и `declarations`; не помещайте его напрямую в `TestBed.imports`.
- Получайте backend URL через `CoreConfigService`. Учитывайте, что Docker `API_URL` сейчас подставляется во время image build.
- Не используйте UI visibility или client-side role state как проверку доступа.
- Перед включением внешнего URL в `HttpClient` ограничьте JWT interceptor backend origin.
- Не дублируйте server state без правила invalidation. При использовании `BaseService` явно решите, нужны ли cache и pagination.
- Сохраняйте пользовательские фильтры в query params там, где это уже делает страница; не вводите второй несогласованный store.
- Завершайте долгоживущие subscriptions при уничтожении component/service connection.
- Dark theme реализуйте существующими Bootstrap variables и `:host-context`; component subscription добавляйте только для поведения, которое нельзя выразить CSS.

## Проверка

- Backend-only изменение: запустите targeted integration test, затем `mise run test-integration`. Для SQL/search/report изменений дополнительно запустите PostgreSQL и MySQL variants через Docker.
- Frontend-only изменение: запустите `npm run build`. Unit test command фиксируйте отдельно; общий suite сейчас имеет известные compile errors, не связанные с production build.
- Изменение API DTO, auth или WebSocket contract проверяйте с обеих сторон и тестируйте как один вертикальный slice.
- Изменение Docker/env проверяйте отдельно для local development и container deployment: эти пути загружают конфигурацию по-разному.
