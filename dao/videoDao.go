package dao

import (
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/xwxb/MiniDouyin/utils/jsonUtils"
)

type TableVideo struct {
	Id            int64     `gorm:"primary_key,auto_increment" json:"id,omitempty"`
	UserId        int64     `gorm:"column:user_id" json:"-"`
	PlayUrl       string    `gorm:"column:play_url" json:"play_url,omitempty"`
	CoverUrl      string    `gorm:"column:cover_url" json:"cover_url"`
	FavoriteCount int64     `gorm:"column:favorite_count" json:"favorite_count,omitempty"`
	CommentCount  int64     `gorm:"column:comment_count" json:"comment_count,omitempty"`
	Author        TableUser `gorm:"foreignKey:Id;references:UserId" json:"author"`
	IsFavorite    bool      `gorm:"-" json:"is_favorite,omitempty"`
	Title         string    `gorm:"column:title" json:"title,omitempty"`
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`

	UserId     int64     `json:"-"`
	Author     User      `gorm:"foreignKey:UserId" json:"author"` //在User中默认会使用Id作为外键的映射值（即外键引用参考值）
	CreatedAt  time.Time `gorm:"column:create_time"`
	Title      string    `json:"title,omitempty"`
	IsFavorite bool      `gorm:"is_favorite" json:"is_favorite"`
}

func (Video) TableName() string {
	return "video"
}

func (video TableVideo) TableName() string {
	return "video"
}

func (User) TableName() string {
	return "user"
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

	if err := Db.Model(&TableVideo{}).
		Preload("Author").
		Joins("left join user u on user_id = u.id").Where("video.id = ?", id).
		Find(&video).Error; err != nil {
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
//	@Description: 根据userId获取用户发布的视频的信息
//	@param userId
//	@return string 返回的格式是json格式
//	@return error
//	使用了联表查询，将作者信息映射到User里面
func GetPublishVideoInfoListByUserId(userId int64) (string, error) {
	var publicVideo []TableVideo
	err := Db.Model(&TableVideo{}).
		Preload("Author").
		Joins("left join user u on user_id = u.id").Where("user_id = ?", userId).
		Find(&publicVideo).Error

	if err != nil {
		log.Println("failed to get PublishVideoInfoList by userId")
	}
	return jsonUtils.MapToJson(publicVideo), err
}

// GetVideoByCreatedTime
//
//	@Description: 根据时间查询晚于该时间发布的视频信息
//	@param lastTime
//	@return string
//	@return error
func GetVideoByCreatedTime(lastTime time.Time) (string, error) {
	var publicVideo []Video

	// Db.AutoMigrate(&Video{})

	err := Db.Model(&Video{}).
		Preload("Author").
		Order("create_time desc").
		Limit(5).
		Joins("left join user u on user_id = u.id"). //当前视频中能在user表中找到对应人，应该是配合外键
		Where("create_time < ?", lastTime).
		Find(&publicVideo).Error

	if err != nil {
		log.Println("视频查询出现问题")
		return "", err
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

func UpdateFavoriteCountByVideoId(videoId int64) bool {
	if err := Db.Model(&Video{}).
		Where("id = ?", videoId).
		Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).
		Error; err != nil {
		log.Println("更新视频喜欢数量失败")
		return false
	}
	return true
}
