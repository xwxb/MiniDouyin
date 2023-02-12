package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/middleware/cos"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type VideoListResponse struct {
	Response
	// 我发现 Video 和 TableVideo 字段基本相同,所以直接返回了 TableVideo
	VideoList []dao.TableVideo `json:"video_list"`
}

// Publish save upload file to public directory
func Publish(c *gin.Context) {
	data, err := c.FormFile("data")
	if err == nil {
		Path := "./public/" //存储路径
		t := time.Now()
		date := t.Format("20060102")
		pathTmp := Path + "/ " + date + "/" //以当天日期命名存储文件夹
		if isDirExists(pathTmp) {
			log.Println("目录存在")
		} else {
			log.Println("目录不存在")
			err := os.Mkdir(pathTmp, 0777) //创建存储文件夹并设置0777权限
			if err != nil {
				//log.Fatal(err)
				c.JSON(http.StatusOK, Response{
					StatusCode: -1,
					StatusMsg:  "mkdir failed",
				})
				log.Printf("创建目录时出错了:\n%v", err)
				return
			}
		}
		//上传文件重命名
		file_name := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+100000) + path.Ext(data.Filename)
		uperr := c.SaveUploadedFile(data, pathTmp+file_name) //文件另存为…
		if uperr == nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  data.Filename + " uploaded successfully",
			})
			log.Printf("%s 上传成功！！", file_name)
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 2,
				StatusMsg:  "upload failed",
			})
			log.Printf("上传时出错了:\n%v", uperr)
			return
		}
		videoURL, _ := cos.UploadToCOS(date, file_name, pathTmp+file_name)
		log.Printf("视频URL是：%v", videoURL)

	} else {
		c.JSON(200, gin.H{"status": 1, "msg": "上传失败"})
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		log.Printf("上传时出错了:\n%v", err)
		return
	}
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

// 目录是否存在
func isDirExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
