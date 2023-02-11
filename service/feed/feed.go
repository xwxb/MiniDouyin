package feed

import (
	"encoding/json"
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
	jsonErr := json.Unmarshal([]byte(VideoList), &FeedList)
	if jsonErr != nil {
		log.Println("解码错误")
		return nil
	}
	//fmt.Printf("未登入获取的feed流：%v\n", FeedList)
	return FeedList
}

func GetFeedByUserId(latestTime time.Time, userId int64) []dao.TableVideo {
	var FeedList []dao.TableVideo
	VideoList, err := dao.GetVideoByCreatedTime(latestTime)

	jsonErr := json.Unmarshal([]byte(VideoList), &FeedList)
	if jsonErr != nil {
		log.Println("解码错误")
		return nil
	}

	if err != nil {
		log.Printf("err = %v\n", err)
	}

	// 封装 isFavorite 和 isFollow
	for k, v := range FeedList {
		FeedList[k].IsFavorite = dao.JudgeFavorByUserId(userId, v.Id)
		FeedList[k].Author.IsFollow, _ = dao.IsFollowed(userId, v.Author.Id)
	}
	//fmt.Printf("登入获取的feed流：%v\n", FeedList)
	return FeedList
}
