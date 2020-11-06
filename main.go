package main

import (
	"fmt"
	"github.com/degary/CloudRestaurant/controller"
	"github.com/degary/CloudRestaurant/tool"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		panic(err)
	}

	_, err = tool.OrmEngine(cfg)
	if err != nil {
		fmt.Errorf("初始化orm engine错误: %s", err.Error())
		return
	}

	port := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	app := gin.Default()
	registerRouter(app)
	app.Run(port)
}

func registerRouter(router *gin.Engine) {
	new(controller.HelloController).Router(router)
	new(controller.MemberController).Router(router)
}
