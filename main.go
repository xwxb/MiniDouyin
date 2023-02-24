package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/middleware/cos"
	"github.com/xwxb/MiniDouyin/service"
	"github.com/xwxb/MiniDouyin/utils/directoryUtils"

	// //测试相关包
	"log"
    "net/http"
    _ "net/http/pprof"//自动注册到本程序监听的端口
	// "os"
    // "runtime/pprof"
)

func main() {
	// pprof
	go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

	// f, err := os.Create("cpu.pprof")
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer f.Close()

    // err = pprof.StartCPUProfile(f)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer pprof.StopCPUProfile()

	// my application
	go service.RunMessageServer()
	go cos.UploadHandle()
	// 每10秒扫描一次public目录，删除空文件夹
	go directoryUtils.DeleteDir("./public/")

	r := gin.Default()

	initRouter(r)

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
