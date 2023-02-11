package controller

import (
	"net/http"
	"strconv"
	"strings"

	"time"
	"fmt"
	_"log"

	"github.com/gin-gonic/gin"

	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/service/feed"
	"github.com/xwxb/MiniDouyin/module"
)

type FeedResponse struct {
	Response
	VideoList []dao.Video `json:"video_list,omitempty"`
	NextTime  int64            `json:"next_time,omitempty"`
}


// Feed same demo video list for every request
func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	auth := c.Query("token")
	//fmt.Println("用户token" + auth)
	fmt.Println("传入的时间:" + inputTime)

	t, _ := strconv.ParseInt(inputTime, 10, 64)
	lastTime := time.Unix(t/1000, 0)

	var videos []dao.Video
	var feedErr error
	if auth == "" {
		videos, feedErr = feed.GetFeed(lastTime)// ？先这样，感觉不是很理解这个用上次获取feed时间的逻辑
	} else {
		auth = strings.Fields(auth)[1]
		user, _ := module.JwtParseUser(auth) // 从 token 解析出 user
		videos, feedErr = feed.GetFeedByUserId(lastTime, user.Id)
	}

	if feedErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "视频查询出现问题"})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "success"},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})

}
