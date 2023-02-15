package dao

import (
	"gorm.io/gorm"
	"log"
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

//返回值来指示是否重复操作
func UpFavor(userId int64, videoId int64) (bool, error) {
	fav := &TableFavor{UserId: userId, VideoId: videoId}

	//如果软删除过了，就执行里面的处理，否则直接创建这条记录; 测试创建成功
	if found := (Db.Unscoped().Where(&fav).First(&fav).Error == nil); found {
		if fav.DeletedAt.Valid {//如果有软删除记录，那么不用重新创建
			// If "DeletedAt.Valid" is true, it's deleted.
			fav.DeletedAt.Valid = false
			// Increase the value of "favorite_count" by 1
			Db.Model(&Video{}).
			Where("id = ?", videoId).
			Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
			return false, nil
		} else {//软删除记录无效的情况（其实我也不知道什么时候会有这种情况）
			return false, nil
		}
	}

	//判断是否已经存在这条记录，没找到才正常添加
	if FavFoundErr := Db.Where(&fav).First(&fav).Error; FavFoundErr != nil {
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
	return true, nil 
}

func UnFav(userId int64, videoId int64) (bool, error) {
	fav := &TableFavor{UserId: userId, VideoId: videoId}

	//如果软删除过了，就执行里面的处理，否则直接创建这条记录; 测试创建成功
	if found := (Db.Unscoped().Where(&fav).First(&fav).Error == nil); found {
		if fav.DeletedAt.Valid  {//有软删除记录，说明重复操作了
			return true, nil
		} else {//无效的软删除记录
			return false, nil
		}
	}

	//不然就是没删除过，直接软删除相应关系
	err := Db.Where(&fav).Delete(&fav).Error
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	//数据库视频表点赞数 - 1
	Db.Model(&Video{}).
	Where("id = ?", videoId).
	Update("favorite_count", gorm.Expr("favorite_count - ?", 1))

	return false, nil
}