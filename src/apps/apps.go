package apps

import (
	"context"
	"fmt"
	"net/http"
	"socious-id/src/apps/utils"
	"socious-id/src/apps/views"
	"socious-id/src/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
)

func Init() *gin.Engine {

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		c.Set("ctx", ctx)
		c.Next()
	})

	//Uploader
	uploader := &utils.GCSUploader{
		CDNUrl:          config.Config.Upload.CDN,
		BucketName:      config.Config.Upload.Bucket,
		CredentialsFile: config.Config.Upload.Credentials,
	}
	router.Use(func(c *gin.Context) {
		c.Set("uploader", uploader)
		c.Next()
	})

	//Cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config.Cors.Origins,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte(config.Config.Secret))
	router.Use(sessions.Sessions("socious-id-session", store))

	//Cache
	router.Use(views.NoCache())

	if config.Config.Debug {
		router.Static("/statics", config.Config.Statics)
	}

	router.LoadHTMLGlob(fmt.Sprintf("%s/*.html", config.Config.Templates))

	views.Init(router)

	//docs
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	router.GET("/docs", gin.WrapH(middleware.SwaggerUI(opts, nil)))
	router.GET("/swagger.yaml", gin.WrapH(http.FileServer(http.Dir("./docs"))))

	return router
}

func Serve() {
	router := Init()
	router.Run(fmt.Sprintf("0.0.0.0:%d", config.Config.Port))
}
