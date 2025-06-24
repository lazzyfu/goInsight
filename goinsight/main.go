package main

import (
	"embed"
	"flag"
	"fmt"
	"goInsight/bootstrap"
	"goInsight/global"
	"goInsight/middleware"
	"goInsight/routers"
	"io/fs"
	"net/http"
	"os"
	"strings"

	commonRouter "goInsight/internal/common/routers"
	dasRouter "goInsight/internal/das/routers"
	inspectRouter "goInsight/internal/inspect/routers"
	ordersRouter "goInsight/internal/orders/routers"
	userRouter "goInsight/internal/users/routers"

	"github.com/gin-gonic/gin"
)

// Define version
var version string

// Read local config file
var configFile = flag.String("config", "config.yaml", "config file")

//go:embed dist
var staticFS embed.FS

const mediaDir = "./media"

func setupStaticFiles(r *gin.Engine) error {
	// Embedded file system
	st, err := fs.Sub(staticFS, "dist")
	if err != nil {
		return fmt.Errorf("Error accessing embedded filesystem: %w", err)
	}
	r.StaticFS("/static", http.FS(st))

	// Provide other non-embedded file system
	if _, err := os.Stat(mediaDir); os.IsNotExist(err) {
		if err := os.MkdirAll(mediaDir, os.ModePerm); err != nil {
			return fmt.Errorf("Failed to create media directory: %w", err)
		}
	}
	r.Static("/media", mediaDir)

	// Default avatar file
	r.StaticFile("/avatar2.jpg", "dist/avatar2.jpg")
	return nil
}

func setupNoRoute(r *gin.Engine) {
	// Fix 404 issue on page refresh
	r.NoRoute(func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
			if content, err := staticFS.ReadFile("dist/index.html"); err == nil {
				c.Header("Content-Type", "text/html; charset=utf-8")
				c.Data(http.StatusOK, "text/html; charset=utf-8", content)
				return
			}
		}
		c.String(http.StatusNotFound, "Not Found")
	})
}

func setupRootRoute(r *gin.Engine) {
	// Root route
	r.GET("/", func(c *gin.Context) {
		if data, err := staticFS.ReadFile("dist/index.html"); err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func RunServer() {
	// Production mode
	if global.App.Config.App.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize authentication middleware
	var err error
	if global.App.JWT, err = middleware.InitAuthMiddleware(); err != nil {
		fmt.Println("Failed to initialize authentication middleware:", err)
		return
	}

	// Load route configs for multiple APPs
	routers.Include(
		userRouter.Routers,
		commonRouter.Routers,
		inspectRouter.Routers,
		dasRouter.Routers,
		ordersRouter.Routers,
	)

	// Initialize router
	r := routers.Init()

	// Static files and routes
	if err := setupStaticFiles(r); err != nil {
		fmt.Println(err)
		return
	}
	setupNoRoute(r)
	setupRootRoute(r)

	// Error handling
	r.Use(gin.Recovery())

	// Start server
	if err := r.Run(global.App.Config.App.ListenAddress); err != nil {
		fmt.Println("Failed to start server: ", err.Error())
	}
}

func main() {
	if version != "" {
		fmt.Println("goInsight Version:", version)
	}
	flag.Parse()
	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		fmt.Printf("Config file %s does not exist, you can also specify the config file path with -config parameter\n", *configFile)
		os.Exit(1)
	}
	bootstrap.InitializeConfig(*configFile)
	bootstrap.InitializeLog()
	global.App.DB = bootstrap.InitializeDB()
	global.App.Redis = bootstrap.InitializeRedis()
	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()
	bootstrap.InitializeCron()
	RunServer()
}
