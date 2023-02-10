package main

import (
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/service"

	"github.com/gin-gonic/gin"
)

func main() {
	dao.Init()
	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
