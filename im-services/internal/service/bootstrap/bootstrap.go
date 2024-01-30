package bootstrap

import (
	"im-services/internal/api/services"
	"im-services/internal/config"
	"im-services/internal/middleware"
	"im-services/internal/service/client"
	coroutine_pool "im-services/pkg/coroutine_pool"
	"im-services/pkg/logger"
	"im-services/pkg/model"
	"im-services/pkg/nsq"
	"im-services/pkg/redis"
	router2 "im-services/pkg/router"
	"im-services/server"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

// Start 启动服务方法
func Start() {

	r := gin.Default()

	go client.ImManager.Start()

	r.Use(middleware.Recover)

	setRoute(r)

	gin.SetMode(config.Conf.Server.Mode)

	server.StartGrpc()

	_ = r.Run(config.Conf.Server.Listen)
}

// LoadConfiguration 加载连接池
func LoadConfiguration() {

	setUpLogger()

	model.InitDb()
	redis.InitClient()

	coroutine_pool.ConnectPool()

	_ = nsq.InitNewProducerPool()
	// todo 消费逻辑可以单独抽离

	services.InitChatBot()

}

// 初始化日志方法
func setUpLogger() {
	logger.InitLogger(
		config.Conf.Log.FileName,
		config.Conf.Log.MaxSize,
		config.Conf.Log.MaxBackup,
		config.Conf.Log.MaxAge,
		config.Conf.Log.Compress,
		config.Conf.Log.Type,
		config.Conf.Log.Level,
	)
}

// 注册路由方法
func setRoute(r *gin.Engine) {
	router2.RegisterApiRoutes(r)
	router2.RegisterWsRouters(r)
}
