# Notification service

Notification service for collecting notifications and sending their notifications to clients.

Conventionally, there is automation or services that must transmit notifications to the user.

Users register and collect notifications for themselves.
On Linux, the client displays notifications using standard means.


## Installation and configuration

1. clone repository
2. in your Postgres create a database and user
3. install goose
4. roll up migrations

```shell
git clone https://github.com/decole/notify-server.git
```

```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Roll up migrations

```shell
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgresql://goose:password@127.0.0.1:8092/go_migrations?sslmode=disable

goose -dir ./migrations up

# or

goose -dir ./migrations postgres "postgresql://goose:password@127.0.0.1:8092/go_migrations?sslmode=disable" up
```

## How to make a service

```shell
go build

# see in folder executable go binary file - 'notify-server'

sudo nano /lib/systemd/system/notify-service.service
```

```shell
[Unit]
Description=go notification service
[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/home/<your-user-name>/notify-service
ExecStart=/home/<your-user-name>/notify-service/notify-server
[Install]
WantedBy=multi-user.target
```


# RTFM

Migrations

https://habr.com/ru/articles/780280/

```shell
# creating migrations
goose -dir ./migrations create init.sql
```

