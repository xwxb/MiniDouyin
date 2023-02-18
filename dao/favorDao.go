package dao

import (
	"gorm.io/gorm"
	"log"
	"errors"
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

// 返回值来指示是否重复操作
func UpFavor(userId int64, videoId int64) (bool, error) {
	fav := &TableFavor{UserId: userId, VideoId: videoId}

	//包括软删除一起查找是否存在这条记录
	if found := (Db.Unscoped().Where(&fav).First(&fav).Error == nil); found {
		if fav.DeletedAt.Valid { //如果有软删除记录，那么不用重新创建
			// If "DeletedAt.Valid" is true, it's deleted.
			// fav.DeletedAt.Valid = false
			Db.Model(&fav).Unscoped().Where(&fav).Update("deleted_at", nil)
			log.Println("将软删除设置为了无效")
			
			// Increase the value of "favorite_count" by 1
			Db.Model(&Video{}).
				Where("id = ?", videoId).
				Update("favorite_count", gorm.Expr("favorite_count + ?", 1))

			return false, nil
		} else { //否则说明重复点赞
			log.Println("检测到重复点赞")
			return true, errors.New("repeat operation")
		}
	}

	//没有这条记录，正常执行点赞操作
	if err := Db.Save(&fav).Error; err != nil {
		log.Println(err.Error())
		return false, err
	}
	//数据库视频表点赞数 + 1
	Db.Model(&Video{}).
		Where("id = ?", videoId).
		Update("favorite_count", gorm.Expr("favorite_count + ?", 1))

	return false, nil
}

func UnFav(userId int64, videoId int64) (bool, error) {
	// log.Println("执行软删除操作")
	fav := &TableFavor{UserId: userId, VideoId: videoId}

	//如果软删除过了，就执行里面的处理，否则直接创建这条记录; 测试创建成功
	if found := (Db.Unscoped().First(&fav).Error == nil); found {
		if fav.DeletedAt.Valid { //有软删除记录，说明重复操作了
			return true, errors.New("repeat operation")
		} else { //有，但没有软删除过，正常删除
			err := Db.Where(&fav).Delete(&fav).Error
			if err != nil {
				log.Println(err.Error())
				log.Println("软删除失败")
				return false, err
			}

			//数据库视频表点赞数 - 1
			Db.Model(&Video{}).
				Where("id = ?", videoId).
				Update("favorite_count", gorm.Expr("favorite_count - ?", 1))

			return false, nil
		}
	} 

	//不然就是重复操作
	return true, errors.New("repeat operation")
	
}
