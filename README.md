Проект нотификационного сервиса.

Нотификационный сервис для сбора нотификаций и отдача клиентам их нотификаций.

Условно, есть автоматизация или сервисы, которые должны передавать нотификации пользователю.

Пользователи регистрируются и собирают себе нотификации.
В Linux клиент показывает уведомление нотификации стандартными средствами.



Миграции
https://habr.com/ru/articles/780280/


```shell
go install github.com/pressly/goose/v3/cmd/goose@latest

goose -dir ./migrations create init.sql
```

Накатить миграию

```shell
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgresql://goose:password@127.0.0.1:8092/go_migrations?sslmode=disable

goose -dir ./migrations up

# or

goose -dir ./migrations postgres "postgresql://goose:password@127.0.0.1:8092/go_migrations?sslmode=disable" up
```

