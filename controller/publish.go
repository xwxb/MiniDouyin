package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	Response
	// 我发现 Video 和 TableVideo 字段基本相同,所以直接返回了 TableVideo
	VideoList []dao.TableVideo `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		log.Println("获取当前用户id失败!")
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取失败",
			},
		})
	}

	id, _ := strconv.ParseInt(userId, 10, 64)
	var publicVideoInfo []dao.TableVideo

	stringInfo, err := dao.GetPublishVideoInfoListByUserId(id)
	if err != nil {
		fmt.Printf("获取用户列表失败:%v\n", err)
	}

	jsonErr := json.Unmarshal([]byte(stringInfo), &publicVideoInfo)
	if jsonErr != nil {
		fmt.Println("解码失败")
	}

	//fmt.Printf("获取到的列表为:"+"\n"+"%v\n", stringInfo)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: publicVideoInfo,
	})
}
