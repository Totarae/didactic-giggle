# cmd/gophermart

В данной директории будет содержаться код накопительной системы лояльности, который скомпилируется в бинарное
приложение.

Frameworks:
zap - логирование
chi - роутер
migrations - миграции БД

```
mockgen -source="cmd/internal/handlers/auth.go" -destination="cmd/internal/handlers/mocks/auth_mock.go" -package=mocks
```
