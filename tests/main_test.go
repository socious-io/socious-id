package tests_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"socious-id/src/apps"
	"socious-id/src/config"
	"strings"
	"testing"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	configPath   string
	router       *gin.Engine
	db           *sqlx.DB
	focused      = false
	authExecuted = false
)

// Setup the test environment before any tests run
var _ = BeforeSuite(func() {
	db, router = setupTestEnvironment()
})

// Drop the database after all tests have run
var _ = AfterSuite(func() {
	teardownTestEnvironment(db)
})

func TestSuite(t *testing.T) {
	checkFocus()
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("Socious Test Suite", Ordered, func() {
	Context("Auth", authGroup)
})

func init() {
	// We back to root dir on execute tests
	os.Chdir("../")
	// Define a flag for the config path
	flag.StringVar(&configPath, "c", "test.config.yml", "Path to the configuration file")
}

func replaceAny(a, b gin.H) {
	for key, valueA := range a {
		if valueB, exists := b[key]; exists {
			// If the value in a is a map, recurse into it
			if mapA, ok := valueA.(gin.H); ok {
				if mapB, ok := valueB.(gin.H); ok {
					replaceAny(mapA, mapB)
				}
			} else if valueA == "<ANY>" {
				// Replace "<ANY>" with the corresponding value from b
				a[key] = valueB
			}
		}
	}
}

func decodeBody(responseBody io.Reader) gin.H {
	body := gin.H{}
	decoder := json.NewDecoder(responseBody)
	decoder.Decode(&body)
	return body
}
func bodyExpect(body, expect gin.H) {
	replaceAny(expect, body)
	Expect(body).To(Equal(expect))
}

func setupTestEnvironment() (*sqlx.DB, *gin.Engine) {
	config.Init(configPath)
	parsedURL, _ := url.Parse(config.Config.Database.URL)
	db := database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})

	schemaFile := fmt.Sprintf("%s/schema.sql", config.Config.Database.SqlDir) // Adjust the path if needed
	schemaContent, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatalf("Failed to read schema.sql: %v", err)
	}
	query := strings.ReplaceAll(string(schemaContent), "$1", parsedURL.User.Username())
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to execute schema.sql: %v", err)
	}
	log.Println("Schema applied successfully!")

	m, err := migrate.New(
		fmt.Sprintf("file://%s", config.Config.Database.Migrations),
		config.Config.Database.URL,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Migrations applied successfully!")
	router := apps.Init()

	return db, router
}

func teardownTestEnvironment(db *sqlx.DB) {
	db.Close()
	if err := database.DropDatabase(config.Config.Database.URL); err != nil {
		log.Fatalf("Dropping database %v", err)
	}
}

func checkFocus() {
	for _, arg := range os.Args[1:] {
		if strings.Contains(arg, "focus") {
			focused = true
			break
		}
	}
}
