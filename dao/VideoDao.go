package dao

import "log"

type TableVideo struct {
	Id            int64  `gorm:"primary_key,auto_increment"`
	UserId        int64  `gorm:"column:user_id"`
	PlayUrl       string `gorm:"column:play_url"`
	CoverUrl      string `gorm:"column:cover_url"`
	FavoriteCount int64  `gorm:"column:favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count"`
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
