package cos

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/xwxb/MiniDouyin/config"
	"github.com/xwxb/MiniDouyin/dao"
	"log"
	"net/http"
	"net/url"
	"time"
)

// 上传视频队列所需的参数
type UploadStruct struct {
	Date string `json:"date"`
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
	User *dao.TableUser `json:"user"`
	Title string `json:"title"`
}
var UploadChan = make(chan UploadStruct, 10)	//缓冲区设为10

func GetClient() *cos.Client {
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶 region 可以在 COS 控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(config.COSURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过配置文件获取密钥
			SecretID: config.COSSecretID,
			SecretKey: config.COSSecretKey,
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			//SecretID: os.Getenv("SECRETID"), // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			//SecretKey: os.Getenv("SECRETKEY"), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	return client

}

// 上传到COS并返回视频URL
func UploadToCOS(date string, filename string, filepath string) (videoURL string, err error) {
	key := date + "/" + filename
	videoDetails, _, err := GetClient().Object.Upload(context.Background(), key, filepath, nil)
	if err != nil {
		log.Printf("上传COS时发生错误：%v", err)
		return "", err
	}
	videoURL = videoDetails.Location
	return videoURL, err
}

// 项目启动时调用，监听并获取channel中的数据实现上传COS操作
func UploadHandle() {
	// recover
	defer func() {
		if errRecover := recover(); errRecover != nil {
			// TODO 记录日志等等，我也不知道这是干嘛的，所以省略了
			fmt.Println("errRecover: ", errRecover)
		}
	}()
	// 轮询、监听channel中的数据，如果收到了就构建请求
	for UploadData := range UploadChan{
		fmt.Println("从channel中轮询接受数据并发送请求")
		// 获取数据结构
		currDate := UploadData.Date
		currFilename := UploadData.Filename
		currFilepath := UploadData.Filepath
		userID := UploadData.User.Id
		title := UploadData.Title
		videoURL, _ := UploadToCOS(currDate, currFilename, currFilepath)
		log.Printf("视频URL是：%v", videoURL)
		// 将数据写入数据库
		video := dao.TableVideo{
			UserId: userID,
			PlayUrl: videoURL,
			CoverUrl: videoURL,
			FavoriteCount: 0,
			CommentCount: 0,
			Title: title,
			CreateTime: time.Now(),
		}
		resp := dao.CreatePublishVideo(&video)
		if resp {
			log.Println("视频数据已写入数据库")
		}
	}
}

// 异步往channel中发送数据（生产者）
func ReportDataToUploadChannel(upload UploadStruct) {
	// recover
	defer func() {
		if errRecover := recover(); errRecover != nil {
			// TODO 记录日志等等，我也不知道这是干嘛的，所以省略了
			fmt.Println("errRecover: ", errRecover)
		}
	}()
	select {
	// 只往channel中发送数据
	case UploadChan <- upload:
	// 缓冲区满了记录一下
	default:
		fmt.Println("缓冲区满了...")
		// TODO 记录log等等
	}
}
