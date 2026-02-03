package main

import (
	"moon/internal/repository"
	"moon/internal/repository/dao"
	"moon/internal/service"
	"moon/internal/web"
	"moon/internal/web/jwt"
	"moon/internal/web/middleware"
	"moon/ioc"
	"moon/pkg/ginx"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件失败: " + err.Error())
	}

	log := ioc.InitLogger()
	ginx.L = log
	log.Info("应用启动中...")

	db := ioc.InitDB()
	rdb := ioc.InitRedis()
	log.Info("数据库和Redis连接成功")

	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewGORMUserRepository(userDAO)
	userService := service.NewUserService(userRepo)

	jwtHdl := jwt.NewRedisJWTHandler(rdb)
	userHandler := web.NewUserHandler(userService, jwtHdl)

	gin.SetMode(viper.GetString("gin.mode"))

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())

	log.Info("路由配置完成")

	jwtMiddleware := middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).CheckLogin()
	router.Use(jwtMiddleware)

	router.GET("/health", func(ctx *gin.Context) {
		log.Info("收到健康检查请求")
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	userHandler.RegisterRoutes(router)

	server := &ginx.Server{
		Addr:   viper.GetString("server.addr"),
		Engine: router,
	}

	log.Info("服务启动在 " + server.Addr)
	if err := server.Start(); err != nil {
		panic("服务启动失败: " + err.Error())
	}
}
