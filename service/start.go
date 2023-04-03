package service

import (
	"fmt"
	"github.com/FelixWuu/chatgpt_self_use/config"
	"github.com/FelixWuu/chatgpt_self_use/routes"
	"github.com/FelixWuu/chatgpt_self_use/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var router *gin.Engine
var once sync.Once

func StartWebService() {
	setupRoute()
	initTemplateDir()
	initStaticServer()
	cfg := config.Inst()

	port := cfg.Address.Port
	listen := cfg.Address.Listen
	err := router.Run(fmt.Sprintf("%s:%d", listen, port))
	if err != nil {
		logger.Errorf("run service failed, error: %v", err)
	}
}

// setupRoute 启动路由
func setupRoute() {
	once.Do(func() {
		router = gin.Default()
		routes.RegisterRoutes(router)
	})
}

// initTemplateDir 初始化 HTML 模板加载路径
func initTemplateDir() {
	router.LoadHTMLGlob("resources/view/*")
}

// initStaticServer 初始化静态文件
func initStaticServer() {
	router.StaticFS("/assets", http.Dir("static/assets"))
	router.StaticFile("logo192.png", "static/logo192.png")
	router.StaticFile("logo512.png", "static/logo512.png")
	router.StaticFile("favicon.ico", "static/favicon.ico")
	router.StaticFile("manifest.json", "static/manifest.json")
}
