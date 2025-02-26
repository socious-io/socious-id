package main

import (
	"socious-id/src/apps"
	"socious-id/src/config"
	"time"

	"github.com/socious-io/gomail"
	"github.com/socious-io/gomq"
	database "github.com/socious-io/pkg_database"
)

func main() {
	config.Init("config.yml")
	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})

	//Initializing GoMQ Library
	gomq.Setup(gomq.Config{
		Url:        config.Config.Nats.Url,
		Token:      config.Config.Nats.Token,
		ChannelDir: "sociousid",
	})
	gomq.Connect()

	gomail.Setup(gomail.Config{
		WorkerChannel: "email",
		MessageQueue:  gomq.Mq,
	})

	apps.Serve()
}
