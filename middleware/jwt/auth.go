package jwt

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/module"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Query("token")
		fmt.Println("auth=", auth)
		if len(auth) == 0 {
			context.Abort()
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		//auth = strings.Fields(auth)[1]
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
			context.Next()
		}
	}
}
