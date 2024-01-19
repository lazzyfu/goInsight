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

	commonRouter "goInsight/internal/app/common/router"
	dasRouter "goInsight/internal/app/das/router"
	ordersRouter "goInsight/internal/app/orders/router"
	userRouter "goInsight/internal/app/users/router"

	"github.com/gin-gonic/gin"
)

var version string

var ConfigFile = flag.String("config", "config.yaml", "config file")

//go:embed dist
var staticFS embed.FS

func RunServer() {
	// 生产环境模式
	if global.App.Config.App.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	global.App.JWT, _ = middleware.InitAuthMiddleware()
	// 加载多个APP的路由配置
	routers.Include(userRouter.Routers)
	routers.Include(commonRouter.Routers)
	routers.Include(dasRouter.Routers)
	routers.Include(ordersRouter.Routers)
	// 初始化路由
	r := routers.Init()
	// 静态资源处理
	st, _ := fs.Sub(staticFS, "dist")
	r.StaticFS("/static", http.FS(st))
	r.StaticFS("/media", http.Dir("media"))
	r.StaticFile("/avatar2.jpg", "dist/avatar2.jpg")
	// 解决页面刷新404的问题
	r.NoRoute(func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			if content, err := os.ReadFile("dist/index.html"); err == nil {
				c.Header("Accept", "text/html")
				c.Data(http.StatusOK, "text/html; charset=utf-8", content)
				return
			}
		}

		c.Writer.WriteHeader(404)
		_, _ = c.Writer.WriteString("Not Found")
	})
	// 根路由
	r.GET("/", func(c *gin.Context) {
		data, err := staticFS.ReadFile("dist/index.html")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
	// 错误处理
	r.Use(gin.Recovery())
	// 启动
	if err := r.Run(global.App.Config.App.ListenAddress); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	if version != "" {
		fmt.Println("Version:", version)
	}

	// 解析输入
	flag.Parse()
	// 初始化配置
	bootstrap.InitializeConfig(*ConfigFile)
	// 初始化日志
	bootstrap.InitializeLog()
	// 初始化数据库
	global.App.DB = bootstrap.InitializeDB()
	global.App.Redis = bootstrap.InitializeRedis()
	// 程序关闭前，释放数据库连接
	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()
	// 初始化cron
	bootstrap.InitializeCron()
	// 启动服务器
	RunServer()
}
