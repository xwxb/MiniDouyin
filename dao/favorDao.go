package dao

import (
	"gorm.io/gorm"
)

type TableFavor struct {
	Id      int64 `gorm:"column:id"`
	UserId  int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
	gorm.DeletedAt
}

func (favor TableFavor) TableName() string {
	return "favor"
}

func JudgeFavorByUserId(userId int64, videoId int64) bool {
	var favor = TableFavor{}
	if err := Db.Where("user_id = ? AND video_id = ?", userId, videoId).Find(&favor).Error; err == nil {
		// 因为使用Find查不到数据会返回空结构体而不会报错，只能用主键判断这条数据是否存在。
		// 主键为0表示不存在
		if favor.Id == 0 {
			return false
		}
		return true
	}
	return false
}
