# AGENTS.md

## Обязательные принципы

- **Не сохраняйте обратную совместимость.** Удаляйте устаревшие пути вместо добавления слоёв совместимости, fallback-механизмов или миграций.
- **Выбирайте самую простую реализацию**, которая полностью удовлетворяет текущим требованиям. Избегайте преждевременных абстракций, избыточной конфигурации и лишних уровней косвенности.
- **Развивайте систему постепенно, слоями.** Сначала доводите минимальную версию до рабочего состояния от начала до конца, затем добавляйте возможности поверх неё. Не жертвуйте рабочим решением ради незавершённой сложности.
- **Сохраняйте модульность компонентов** и чёткое разделение ответственности.
- **Предпочитайте проверенные и хорошо поддерживаемые библиотеки**, если они уменьшают общую сложность или повышают надёжность. Не реализуйте распространённую функциональность заново без веской причины.
- **Максимально используйте существующие зависимости** до написания собственной реализации или добавления пакета. Сначала проверяйте документацию и типы библиотеки.
- **Принимайте архитектурные решения на долгий срок.** Не добавляйте заведомо временные решения, которые потребуется заменить.

## Документация

- [Архитектура системы](.docs/architecture.md)
- [Архитектура backend](.docs/arch-backend.md)
- [Архитектура frontend](.docs/arch-frontend.md)
- [Правила кода и слоёв](.docs/rules.md)

При изменении границ модулей, потоков данных, wiring, инфраструктуры или обязательных команд обновляйте соответствующий документ в том же изменении.

## Рабочие команды

- Версии инструментов задаёт `mise.toml`: Go 1.26, Node 22, `golangci-lint` 2.10.1.
- Backend: `mise run dev-backend`; frontend: `mise run dev-frontend`. Эти задачи загружают корневой `.env` через mise.
- Backend integration, SQLite: `mise run test-integration`.
- Backend integration, PostgreSQL/MySQL: `mise run test-integration-postgresql` / `mise run test-integration-mysql`; обе команды требуют Docker для Testcontainers.
- Один backend-тест: `cd backend && DB_TYPE=sqlite LOG_LEVEL=1 go test -v -count=1 -timeout 300s ./tests/integration/... -run '^TestGetUsers$'`.
- Frontend: из `frontend/board` запускайте `npm ci`, `npm run build`, `npm test -- --watch=false --browsers=ChromeHeadless`.
- Один frontend spec задаётся через `--include='src/path/file.spec.ts'`, но TypeScript всё равно компилирует все spec-файлы. На 2026-08-03 suite падает до запуска тестов из-за удалённого Angular API `async` в трёх spec-файлах `src/app/modules/ui/components/{select,select-model,select-model-multiple}/`.
- В проекте нет настроенных frontend lint/e2e-команд. Не придумывайте их в отчётах о проверке.

## Источники истины

- Считайте `mise.toml`, `backend/go.mod`, `frontend/board/package-lock.json`, `frontend/board/angular.json` и исходный код актуальнее README и CI. GitHub workflows используют старые Go 1.21/Node 18 и не запускают реальные backend integration/frontend unit tests.
- Backend сам не читает `.env`: переменные предоставляет mise, Docker Compose или окружение процесса.
- Не коммитьте корневой `.env`, локальные `frontend/board/{node_modules,dist,.angular}`, backend debug binaries или данные `.tmp`/`.repowise`.
