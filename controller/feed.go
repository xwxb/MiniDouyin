package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/service/fedd"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
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
	log.Println("传入的时间:" + inputTime)

	var lastTime time.Time
	if inputTime != "" {
		latestTime, timeErr := strconv.ParseInt(inputTime, 10, 64)
		lastTime = time.Unix(latestTime, 0)
		if timeErr != nil {
			lastTime = time.Now()
		}
	} else {
		lastTime = time.Now()
	}

	log.Printf("最后的投稿时间: %v\n", lastTime)

	// 未登入直接返回 feed
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: fedd.GetFeed(lastTime),
		NextTime:  time.Now().Unix(),
	})

	// resp, err := getFeed()

	// if err != nil {
	// 	panic("failed to get feed!")
	// }

	// c.JSON(http.StatusOK, resp)
}

// func getFeed() (FeedResponse, error) {
// 	// 数据库连接
// 	dsn := "${root:114514@tcp(47.94.10.223:3306)/mdy?charset=utf8mb4&parseTime=True"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic("failed to connect database")
// 		return nil, err
// 	}

// 	db.AutoMigrate(&Video{})

// 	// todo1 添加逻辑：用户投稿

// 	// 查询逻辑
// 	// todo2 判断用户是否登录

// 	// 如果未登录，根据视频id查询作者信息

// 	//组装返回

// 	var videos []Video

// 	if err := db.Preload("Author").Find(&videos).Error; err != nil {
// 		return nil, err
// 	}

// 	resp := FeedResponse{
// 		Response:  Response{StatusCode: 0},
// 		VideoList: []Video{videos},
// 		NextTime:  time.Now().Unix(),
// 	}

// 	return resp, nil
// }
