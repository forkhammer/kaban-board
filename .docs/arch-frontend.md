# Архитектура frontend

## Bootstrap и маршрутизация

`src/main.ts` загружает `AppModule`; приложение остаётся преимущественно NgModule-based. `app.module.ts` подключает `CoreModule.forRoot`, `UiModule`, `KanbanModule`, Bootstrap wrappers, forms, Hammer и DI-based JWT interceptor. Страницы reports являются standalone-компонентами и lazy-load через `app-routing.module.ts`.

Маршруты: `/` (доска), `/auth`, `/auth/gitlab`, `/sprints`, `/search`, `/reports` и три страницы `/reports/{burndown,burnup,wip}`. Guards и wildcard/404 route отсутствуют; авторизацию обеспечивает backend, а не router.

## Модули

| Каталог | Ответственность |
| --- | --- |
| `app/components` | Route-level shell pages: board, auth/OAuth callback, sprints, search и reports catalog |
| `app/modules/core` | API config, auth/JWT, generic `BaseService`, collection cache, errors, toast и общие модели |
| `app/modules/ui` | Переиспользуемые selectors, form feedback, spinner, account UI, pipes/directives |
| `app/modules/bootstrap-ui` | Централизованные exports ng-bootstrap |
| `app/modules/kanban` | Board, issue/sprint/admin UI, domain models, REST services, modal facades и WebSocket client |
| `app/modules/reports` | Standalone Plotly pages, report DTO и `ReportService` |

`CoreModule.forRoot()` должен импортироваться только в `AppModule`: constructor модуля намеренно падает при повторном импорте.

## Данные и состояние

Глобального store нет. Состояние распределено между component fields, forms/query params, root services с `BehaviorSubject` и RxJS streams.

- `BaseService<T>` реализует list/all/get/save/patch/delete, pagination и opt-in in-memory `CollectionCache`; feature services задают свой endpoint и отключают cache/pagination там, где это необходимо.
- URL query params являются сохраняемым состоянием фильтров board, sprints и reports.
- Issues/users периодически перечитываются с интервалом `environment.autoUpdateIssuesMin`.
- Для выбранного sprint `IssueTableComponent` подписывается на `SprintWebsocketService` и применяет binding create/update/delete/order events к локальному состоянию.
- Modal facade services открывают `NgbModal` и возвращают Promise с результатом.

Оптимистические UI-обновления не всегда имеют rollback. При изменении mutation flow явно определяйте восстановление состояния после HTTP conflict/error и последующее WebSocket событие.

## Auth и API

JWT хранится в `localStorage.accessToken`. `JWTInterceptor` добавляет bearer token ко всем HttpClient requests, а не только к backend origin; до подключения внешнего HTTP API ограничьте interceptor, чтобы не отправить токен третьей стороне.

Frontend role checks управляют отображением и drag/drop, но не защищают данные. Все ограничения должны дублироваться route-level middleware backend.

Production API URL читается из `window.env.apiUrl`, загруженного `index.html` из `assets/env.js`, с fallback на `${window.location.origin}/api`. Development environment жёстко использует `http://localhost:8080/api`.

`assets/env.js` не хранится в git: Dockerfile генерирует его из `API_URL` перед `npm run build`. Это build-time configuration. Обычный локальный production build не создаёт файл, а nginx не proxy-ирует `/api`; deployment обязан предоставить `env.js` или внешний same-origin proxy.

WebSocket URL получается заменой `http` на `ws` в API URL и добавлением `/sprint/:id/ws?token=...`; изменение backend route или auth contract требует синхронного изменения сервиса.

## UI и стили

Глобальные SCSS находятся в `src/styles.scss` и `src/styles/`. Bootstrap 5 подключается и переопределяется через SCSS; тёмная тема использует `data-bs-theme` и сохраняется под `localStorage.dark-mode`. Для локального dark style используйте существующий `:host-context([data-bs-theme="dark"])` pattern вместо подписки компонента на theme, если runtime logic не нужна.

Основные UI-зависимости уже в проекте: ng-bootstrap, Angular CDK drag/drop, Font Awesome, Plotly, HammerJS и lodash. Перед собственной реализацией проверяйте их API и типы; Plotly использует локальный shim `src/types/plotly.d.ts`.

## Правила развития

- Существующие feature/UI components с `standalone: false` объявляйте и при необходимости экспортируйте из `KanbanModule`/`UiModule`; standalone route component перечисляет собственные imports.
- Общий CRUD стройте на `BaseService`, только если совпадают endpoint semantics, pagination и cache behavior. Не подгоняйте нестандартный flow под generic abstraction.
- API URL получайте через `CoreConfigService`; не добавляйте новые hardcoded backend origins.
- Сохраняемые фильтры синхронизируйте с query params по существующему pattern.
- Для subscriptions, привязанных к component lifecycle, используйте `takeUntilDestroyed` или эквивалентное явное завершение.
- При добавлении binding event синхронно обновляйте backend event type/payload, `SprintWebsocketService` и локальную обработку в issue table.
- Не импортируйте Angular symbols из внутренних hashed paths; используйте публичные entrypoints пакета.

## Сборка и тесты

Из `frontend/board`:

```bash
npm ci
npm run build
npm run start
npm test -- --watch=false --browsers=ChromeHeadless
npm test -- --watch=false --browsers=ChromeHeadless --include='src/app/modules/kanban/pipes/filter-by-labels.pipe.spec.ts'
```

Production output: `dist/board/browser`. Отдельных lint и e2e targets нет.

Karma/Jasmine suite компилирует все `*.spec.ts` даже с `--include`. На 2026-08-03 она падает до запуска из-за `async` import из `@angular/core/testing` в spec-файлах select components. Кроме того, многие scaffold specs неверно помещают non-standalone components в `TestBed.imports`. Не заявляйте unit tests пройденными, пока эти общие compile errors не устранены; production build при этом проходит.
