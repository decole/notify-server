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

## REST API

-----

### Register client

`POST localhost:8881/client/signup`

body:

```json
{
    "client": "decole"
}
```

-----

### Check signup client
`<user>` - your registered user

`GET localhost:8881/client/is-signup/decole`

For example: `GET localhost:8881/client/is-signup/<user>`

-----

### Get notify


`GET localhost:8181/notify/<user>`

For example: `GET localhost:8181/notify/decole`

-----

### Send notify by specific client

`POST localhost:8881/notify`

body:

```json
{
  "client": "decole",
  "message": "test message"
}
```

-----

### Send notify by all clients

`POST localhost:8881/notify`

body:

```json
{
  "message": "test message"
}
```

-----

# RTFM

Migrations

https://habr.com/ru/articles/780280/

```shell
# creating migrations
goose -dir ./migrations create init.sql
```

