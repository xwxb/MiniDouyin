package feed

import (
	"encoding/json"
	"log"
	"time"
	"github.com/xwxb/MiniDouyin/dao"
	
)

func GetFeed(latestTime time.Time) (time.Time, []dao.Video, error) {
	var FeedList []dao.Video

	VideoList, err := dao.GetVideoByCreatedTime(latestTime)
	if err != nil {
		log.Printf("err = %v\n", err)
	}

	
	jsonErr := json.Unmarshal([]byte(VideoList), &FeedList)
	if jsonErr != nil || len(FeedList) == 0 {
		log.Println("解码错误或数据库已空")
		return time.Time{}, nil, jsonErr
	}

	//fmt.Printf("未登入获取的feed流：%v\n", FeedList)
	// log.Printf("当前获取到的视频数%v\n", len(FeedList))

	return FeedList[len(FeedList) - 1].CreatedAt, FeedList, nil
}

func GetFeedByUserId(latestTime time.Time, userId int64) (time.Time, []dao.Video, error) {
	var FeedList []dao.Video
	VideoList, err := dao.GetVideoByCreatedTime(latestTime)

	jsonErr := json.Unmarshal([]byte(VideoList), &FeedList)
	if jsonErr != nil || len(FeedList) == 0 {
		log.Println("解码错误或数据库已空")
		return time.Time{}, nil, jsonErr
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
	return FeedList[len(FeedList) - 1].CreatedAt, FeedList, nil
}