package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/module"
	"github.com/xwxb/MiniDouyin/service/feed"
	_ "log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []dao.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//获取用户传来的参数
	inputTime := c.Query("latest_time")
	auth := c.Query("token")
	//fmt.Println("用户token" + auth)
	// fmt.Println("传入的时间:" + inputTime)

	//时间处理
	t, _ := strconv.ParseInt(inputTime, 10, 64)
	lastTime := time.Unix(t, 0)

	if lastTime.After(time.Now()) { //暂时按第一次是一个大时间考虑
		fmt.Println("第一次")
		lastTime = time.Now()
	} else { //否则因为第一次会传来上次结束获取的时间，一定比现在早
		// 给他上一次获取的最后一个视频的时间
		fmt.Println("不是第一次了")
	}
	fmt.Printf("传入的时间 = %v \n", lastTime)

	var videos []dao.Video
	var feedErr error
	var lastestVideoTime time.Time
	if auth == "" {
		lastestVideoTime, videos, feedErr = feed.GetFeed(lastTime) // ？先这样，感觉不是很理解这个用上次获取feed时间的逻辑
	} else {
		auth = strings.Fields(auth)[1]
		user, _ := module.JwtParseUser(auth) // 从 token 解析出 user
		lastestVideoTime, videos, feedErr = feed.GetFeedByUserId(lastTime, user.Id)
	}

	if feedErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "视频查询出现问题"})
		return
	}
	// fmt.Printf("\n\n%v\n\n", videos)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "success"},
		VideoList: videos,
		NextTime:  lastestVideoTime.Unix(),
	})

}
