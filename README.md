# Socious ID

## Database Query system
**Temp Note**: *would update db query related documents here*

## Migration system
**Temp Note**: *would update migration related documents here*

## Quick start
**should take care of matching config file to related connection such as pg and nats**
**develop runner**:
```
$ cd socious-id
$ curl -sSfL https://goblin.run/github.com/air-verse/air | sh
$ cp .tmp.config.yml config.yml
$ sudo docker-compose up -d
$ go get
$ air
``` 