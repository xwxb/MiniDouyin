package fedd

import (
	"encoding/json"
	"fmt"
	"github.com/xwxb/MiniDouyin/dao"
	"log"
	"time"
)

func GetFeed(latestTime time.Time) []dao.TableVideo {
	var FeedList []dao.TableVideo
	VideoList, err := dao.GetVideoByCreatedTime(latestTime)
	if err != nil {
		log.Printf("err = %v\n", err)
	}
	json.Unmarshal([]byte(VideoList), &FeedList)
	fmt.Printf("未登入获取的feed流：%v\n", FeedList)
	return FeedList
}
