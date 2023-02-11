package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/module"
	"github.com/xwxb/MiniDouyin/service/feed"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []dao.TableVideo `json:"video_list,omitempty"`
	NextTime  int64            `json:"next_time,omitempty"`
}

func (Video) TableName() string {
	return "video"
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	auth := c.Query("token")
	//fmt.Println("用户token" + auth)
	fmt.Println("传入的时间:" + inputTime)

	var lastTime time.Time
	latestTime, timeErr := strconv.ParseInt(inputTime, 10, 64)
	if timeErr != nil {
		log.Printf("time err: %v\n", timeErr)
	}

	lastTime = time.Unix(latestTime, 0)

	if lastTime.After(time.Now()) {
		lastTime = time.Now()
	}

	log.Printf("最后的投稿时间: %v\n", lastTime)

	if auth == "" {
		// 未登入直接返回 feed
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: feed.GetFeed(lastTime),
			NextTime:  time.Now().Unix(),
		})
	} else {
		//登入了返回组装好的 feed
		auth = strings.Fields(auth)[1]
		user, _ := module.JwtParseUser(auth) // 从 token 解析出 user
		//log.Printf("传入的用户Id: %v\n", user.Id)
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: feed.GetFeedByUserId(lastTime, user.Id),
			NextTime:  time.Now().Unix(),
		})
	}
}
