package app

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/lazzyfu/goinsight/api"
	"github.com/lazzyfu/goinsight/internal/bootstrap"
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/middleware"
	"github.com/lazzyfu/goinsight/web"

	commonRouter "github.com/lazzyfu/goinsight/internal/common/routers"
	dasRouter "github.com/lazzyfu/goinsight/internal/das/routers"
	inspectRouter "github.com/lazzyfu/goinsight/internal/inspect/routers"
	ordersRouter "github.com/lazzyfu/goinsight/internal/orders/routers"
	userRouter "github.com/lazzyfu/goinsight/internal/users/routers"

	"github.com/gin-gonic/gin"
)

const mediaDir = "./media"

func setupStaticFiles(r *gin.Engine) error {
	// Embedded file system - 映射整个 dist 目录
	distFS, err := fs.Sub(web.StaticFS, "dist")
	if err != nil {
		return fmt.Errorf("Error accessing embedded filesystem: %w", err)
	}

	// 映射 /static 到 dist 根目录（用于访问 avatar2.jpg 等文件）
	r.StaticFS("/static", http.FS(distFS))

	// 映射 assets 目录
	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		return fmt.Errorf("Error accessing assets directory: %w", err)
	}
	r.StaticFS("/assets", http.FS(assetsFS))

	// Provide other non-embedded file system
	if _, err := os.Stat(mediaDir); os.IsNotExist(err) {
		if err := os.MkdirAll(mediaDir, os.ModePerm); err != nil {
			return fmt.Errorf("Failed to create media directory: %w", err)
		}
	}
	r.Static("/media", mediaDir)

	return nil
}

func setupNoRoute(r *gin.Engine) {
	// Fix 404 issue on page refresh
	r.NoRoute(func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
			if content, err := web.StaticFS.ReadFile("dist/index.html"); err == nil {
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
		if data, err := web.StaticFS.ReadFile("dist/index.html"); err == nil {
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
	api.Include(
		userRouter.Routers,
		commonRouter.Routers,
		inspectRouter.Routers,
		dasRouter.Routers,
		ordersRouter.Routers,
	)

	// Initialize router
	r := api.Init()

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

func Run(configFile string) {
	bootstrap.InitializeConfig(configFile)
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
