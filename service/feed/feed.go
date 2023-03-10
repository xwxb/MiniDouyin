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

	return FeedList[len(FeedList)-1].CreatedAt, FeedList, nil
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

	var vids, uids []int64
	for _, v := range FeedList {
		vids = append(vids, v.Id)
		uids = append(uids, v.Author.Id)
	}

	favStats := dao.JudgeFavorByUserIdMult(userId, vids)
	folStats := dao.IsFollowedMult(userId, uids)

	// 封装 isFavorite 和 isFollow
	for i := 0; i < len(FeedList); i++ {
		FeedList[i].IsFavorite = favStats[i]
		// log.Println(FeedList[k].IsFavorite)
		FeedList[i].Author.IsFollow = folStats[i]
		FeedList[i].CommentCount, _ = dao.GetCommentNum(FeedList[i].Id)
	}
	//fmt.Printf("登入获取的feed流：%v\n", FeedList)
	return FeedList[len(FeedList)-1].CreatedAt, FeedList, nil
}
