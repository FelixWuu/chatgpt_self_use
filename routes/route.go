package routes

import (
	"github.com/FelixWuu/chatgpt_self_use/config"
	"github.com/FelixWuu/chatgpt_self_use/controllers"
	"github.com/FelixWuu/chatgpt_self_use/middlewares"
	"github.com/gin-gonic/gin"
)

var rspCtrl = controllers.NewResponseController()

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.Engine) {
	router.Use(middlewares.Cors())
	cfg := config.Inst()
	if len(cfg.Auth.AuthUser) > 0 {
		router.Use(gin.BasicAuth(gin.Accounts{
			cfg.AuthUser: cfg.Auth.AuthPassword,
		}))
	}

	router.GET("/", rspCtrl.Index)
	router.POST("/completion", rspCtrl.Response)
}
