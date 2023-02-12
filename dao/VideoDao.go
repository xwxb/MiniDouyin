package dao

import (
	"github.com/xwxb/MiniDouyin/utils/jsonUtils"
	"log"
	"time"
)

type TableVideo struct {
	Id            int64     `gorm:"primary_key,auto_increment" json:"id,omitempty"`
	UserId        int64     `gorm:"column:user_id" json:"-"`
	PlayUrl       string    `gorm:"column:play_url" json:"play_url,omitempty"`
	CoverUrl      string    `gorm:"column:cover_url" json:"cover_url,omitempty"`
	FavoriteCount int64     `gorm:"column:favorite_count" json:"favorite_count,omitempty"`
	CommentCount  int64     `gorm:"column:comment_count" json:"comment_count,omitempty"`
	Author        TableUser `gorm:"foreignKey:Id;references:UserId"`
	IsFavorite    bool      `gorm:"-" json:"is_favorite,omitempty"`
	Title         string    `gorm:"column:title" json:"title,omitempty"`
	CreateTime    time.Time	`gorm:"column:create_time" json:"create_time,omitempty"`
}

func (video TableVideo) TableName() string {
	return "video"
}

func GetVideoList() ([]TableVideo, error) {
	var videosList []TableVideo
	if err := Db.Find(&videosList).Error; err != nil {
		log.Println(err.Error())
		return videosList, err
	}
	return videosList, nil
}

func GetVideoByVideoId(id int64) (TableVideo, error) {
	video := TableVideo{}
	if err := Db.Where("id = ?", id).First(&video).Error; err != nil {
		log.Println(err.Error())
		return video, err
	}
	return video, nil
}

func GetVideosListByUserId(userId int64) ([]TableVideo, error) {
	var videosList []TableVideo
	if err := Db.Where("user_id = ?", userId).Find(&videosList).Error; err != nil {
		log.Println(err.Error())
		return videosList, err
	}
	return videosList, nil
}

// GetPublishVideoInfoListByUserId
//
//	 @Description: 根据userId获取用户发布的视频的信息
//	 @param userId
//	 @return string 返回的格式是json格式
//	 @return error
//		使用了联表查询，将作者信息映射到User里面
func GetPublishVideoInfoListByUserId(userId int64) (string, error) {
	var publicVideo []TableVideo
	err := Db.Model(&TableVideo{}).
		Preload("Author").
		Joins("left join user u on user_id = u.id").Where("user_id = ?", userId).
		Find(&publicVideo).Error

	if err != nil {
		log.Println("failed")
	}
	return jsonUtils.MapToJson(publicVideo), err
}

// CreatePublishVideo 上传的视频信息添加到数据库
func CreatePublishVideo(video *TableVideo) bool {
	if err := Db.Create(&video).Error; err != nil {
		log.Println("视频插入到数据库时产生错误：", err)
		return false
	}
	return true
}