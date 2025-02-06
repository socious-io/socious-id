# Socious ID

## Database Query system
**Temp Note**: *would update db query related documents here*

## Migration system
**Temp Note**: *would update migration related documents here*

## Quick start
**should take care of matching config file to related connection such as pg and nats**
```
$ cd socious-id
$ cp .tmp.config.yml config.yml
$ sudo docker-compose up -d
$ go get
$ go run cmd/app/main.go
``` 