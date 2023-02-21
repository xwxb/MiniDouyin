package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/middleware/cos"
	"github.com/xwxb/MiniDouyin/middleware/ffmpeg"
	"github.com/xwxb/MiniDouyin/service/user"
	"github.com/xwxb/MiniDouyin/utils/directoryUtils"
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
	// (通过token验证解析后得到的)user(的引用)
	userP := c.MustGet("authUserObj").(*dao.TableUser)

	// 由于前端返回的是form格式数据，这里要用PostFormValue的方法获取title
	title := c.Request.PostFormValue("title")

	// 前端返回的视频数据
	data, err := c.FormFile("data")
	if err == nil {
		Path := "./public" //存储路径
		t := time.Now()
		date := t.Format("20060102")
		pathTmp := Path + "/" + date + "/" //以当天日期命名存储文件夹，./public/20230215/
		if directoryUtils.IsDirExists(pathTmp) {
			fmt.Println("目录存在")
		} else {
			fmt.Println("目录不存在")
			err := os.Mkdir(pathTmp, 0777) //创建存储文件夹并设置0777权限
			if err != nil {
				//log.Fatal(err)
				c.JSON(http.StatusOK, Response{
					StatusCode: -1,
					StatusMsg:  "mkdir failed",
				})
				fmt.Printf("创建目录时出错了:\n%v", err)
				return
			}
		}
		//生成随机种子
		seed := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+100000)
		//上传文件重命名，1676446898304475.mp4
		file_name := seed + path.Ext(data.Filename)
		//封面图片命名，1676446898304475.jpg
		image_name := seed + ".jpg"
		uperr := c.SaveUploadedFile(data, pathTmp+file_name) //文件另存为…
		if uperr == nil {
			//更新用户投稿数，作品数+1
			err := user.AddWorkCount(userP.Id)
			if err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 2,
					StatusMsg:  "database error",
				})
				fmt.Println("更改数据库失败！")
				if err := os.Remove(pathTmp + file_name); err != nil {
					log.Println("视频删除时出错了：", pathTmp+file_name)
				} else {
					fmt.Println("本地视频已删除")
				}
				return
			}

			//截图
			go func() {
				fferr := ffmpeg.Ffmpeg(date+"/"+file_name, date+"/"+image_name)
				if fferr == nil {
					fmt.Println("封面截取成功！")
				} else {
					fmt.Printf("封面截取失败：%v", fferr)
				}
			}()

			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  data.Filename + " uploaded successfully",
			})
			fmt.Printf("%s 上传成功！！", file_name)
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 2,
				StatusMsg:  "upload failed",
			})
			fmt.Printf("上传时出错了:\n%v", uperr)
			return
		}
		// 将所需参数构件好结构体放入channel，实现异步上传到云端
		// 因为是小项目，我的理解是并发应该不会很大，所以采用队列的方法用单个子线程上传，但如果大项目的话可能也许直接开多线程？
		currUpload := cos.UploadStruct{
			Date:      date,
			Filename:  file_name,
			Filepath:  pathTmp + file_name,
			Imagename: image_name,
			Imagepath: pathTmp + image_name,
			User:      userP,
			Title:     title,
		}
		cos.ReportDataToUploadChannel(currUpload)

	} else {
		c.JSON(200, gin.H{"status": 1, "msg": "上传失败"})
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		fmt.Printf("上传时出错了:\n%v", err)
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
