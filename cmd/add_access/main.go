package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/config"
	"time"

	database "github.com/socious-io/pkg_database"
)

var (
	name        = flag.String("n", "example", "access name")
	description = flag.String("d", "example description", "access description")
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

	defer database.Close()
	secret := utils.RandomString(24)
	clientSecret, _ := auth.HashPassword(secret)

	access := &models.Access{
		Name:         *name,
		Description:  *description,
		ClientID:     utils.RandomString(8),
		ClientSecret: clientSecret,
	}
	ctx := context.Background()

	if err := access.Create(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("one time client secret visibility, take copy please")
	fmt.Println(" ------------------------ ")
	fmt.Printf("Client ID : `%s` \n", access.ClientID)
	fmt.Printf("Client Secret : `%s` \n", secret)
	fmt.Println(" ------------------------ ")
}
