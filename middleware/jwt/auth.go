package jwt

import (
	"github.com/xwxb/MiniDouyin/module"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// jwt 中间件, 每次请求一个路径之前验证一下
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Query("token")

		log.Println("auth=", auth)
		if len(auth) == 0 {
			context.Abort()
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		auth = strings.Fields(auth)[1] // 去掉前面的 “Bearer ”
		usr, err := module.JwtParseUser(auth)
		if err != nil {
			context.Abort()
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
			log.Fatal("err:", err)
		} else {
			context.Set("authUserObj", usr)
			//fmt.Printf("usr = %v\n", usr)
			context.Next()
		}
	}
}
