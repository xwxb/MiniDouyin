package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/middleware/cos"
	"github.com/xwxb/MiniDouyin/service"
	"github.com/xwxb/MiniDouyin/utils/directoryUtils"
)

func main() {
	go service.RunMessageServer()
	go cos.UploadHandle()
	// 每10秒扫描一次public目录，删除空文件夹
	go directoryUtils.DeleteDir("./public/")

	r := gin.Default()

	initRouter(r)

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
