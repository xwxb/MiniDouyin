package controller

import (
	"net/http"
	// "strconv"
	// "strings"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/xwxb/MiniDouyin/config"
	_"github.com/xwxb/MiniDouyin/dao"

	_"fmt"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64            `json:"next_time,omitempty"`
}

func (Video) TableName() string {
	return "video"
}

func (User) TableName() string {
	return "user"
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	// 返回格式参考
	// c.JSON(http.StatusOK, FeedResponse{
	// 	Response:  Response{StatusCode: 0},
	// 	VideoList: DemoVideos,
	// 	NextTime:  time.Now().Unix(),
	// })
	

	
	// 数据库连接
	dsn := config.LoginHead + "?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect database"})
		return
	}
	

	var videos []Video
	if err := db.Preload("Author").Order("create_time desc").Limit(30).Find(&videos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "视频查询出现问题"})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "success"},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})

}
