package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/module"

	"strings"
	"strconv"

	"encoding/json"
	"log"
)

type FavRequset struct {
	Token      string `json:"token,omitempty"`
	VideoId    int64 `json:"video_id,omitempty"`
	ActionType int64  `json:"action_type,omitempty"`
}

type FavRespond struct {
	StatusCode int64 `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	tk := c.Request.URL.Query().Get("token")
	vid := c.Request.URL.Query().Get("video_id")
	at := c.Request.URL.Query().Get("action_type")

	vid_int, err := strconv.ParseInt(vid, 10, 64)
	at_int, err := strconv.ParseInt(at, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	req := FavRequset{
		Token: tk,
		VideoId: vid_int,
		ActionType: at_int,
	}

	// fmt.Printf("\n\n%v\n\n", c.Request)
	// fmt.Printf("\n\n%v\n\n", c.Request.URL)
	// fmt.Printf("token = %v\n", req.Token)
	// fmt.Printf("vid = %v\n", req.VideoId)
	// fmt.Printf("actType = %v\n", req.ActionType)

	//Retrieve the userid
	var uid int64
	if req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 1,
			"status_msg":  "please login first"})
		return
	} else {
		auth := strings.Fields(req.Token)[1]
		user, err := module.JwtParseUser(auth) // 从 token 解析出 user

		if err != nil {
			fmt.Println(err)
		}
		// fmt.Printf("\nuser = %v\n", user)

		uid = user.Id
	}
	// fmt.Printf("\nuid = %v\n", uid)

	resp := FavRespond {
		StatusCode: 0,
	}

	// 考虑是否需要布尔返回值
	if req.ActionType == 1 {
		ifFavorAlready, err :=dao.UpFavor(uid, req.VideoId)

		if err != nil {
			if ifFavorAlready == true {
				resp.StatusMsg = "Repeated like action"
			} else {
				resp.StatusMsg = "like action success!"
			}
		} else {
			resp.StatusMsg = "like action failure!"
		}

	} else if req.ActionType == 2 {
		ifUnFavAlready, err := dao.UnFav(uid, req.VideoId)

		if err != nil {
			if ifUnFavAlready == true {
				resp.StatusMsg = "Repeated unlike action"
			} else {
				resp.StatusMsg = "unlike action success!"
			}
		} else {
			resp.StatusMsg = "unlike action failure!"
		}

		resp.StatusMsg = "dislike action success!"
	} else {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid action type"
	}

	c.JSON(http.StatusOK, resp)
}

// //暂时没发现这个是干嘛的
// func Favor(c *gin.Context) {
// 	token := c.PostForm("token")

// 	if _, exist := usersLoginInfo[token]; !exist {
// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
// 		return
// 	}

// 	data, err := c.FormFile("data")
// 	if err != nil {
// 		c.JSON(http.StatusOK, Response{
// 			StatusCode: 1,
// 			StatusMsg:  err.Error(),
// 		})
// 		return
// 	}

// 	filename := filepath.Base(data.Filename)
// 	user := usersLoginInfo[token]
// 	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
// 	//save file
// 	saveFile := filepath.Join("./favor/", finalName)
// 	if err := c.SaveUploadedFile(data, saveFile); err != nil {
// 		c.JSON(http.StatusOK, Response{
// 			StatusCode: 1,
// 			StatusMsg:  err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, Response{
// 		StatusCode: 0,
// 		StatusMsg:  finalName + " uploaded successfully",
// 	})

// }

// have list by userid
func FavoriteList(c *gin.Context) {
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
	var favorVideoInfo []dao.TableVideo

	stringInfo, err := dao.GetFavorVideoInfoListByUserId(id)
	if err != nil {
		fmt.Printf("获取用户列表失败:%v\n", err)
	}

	jsonErr := json.Unmarshal([]byte(stringInfo), &favorVideoInfo)
	if jsonErr != nil {
		fmt.Println("解码失败")
	}

	//fmt.Printf("获取到的列表为:"+"\n"+"%v\n", stringInfo)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: favorVideoInfo,
	})
}
