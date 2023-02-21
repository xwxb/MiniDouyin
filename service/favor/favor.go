package favor

import (
	"errors"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/service/user"
	"github.com/xwxb/MiniDouyin/utils/jsonUtils"
	"gorm.io/gorm"
	"log"
)

func UpFavor(userId int64, videoId int64) (bool, error) {
	fav := &dao.TableFavor{UserId: userId, VideoId: videoId}
	video, _ := dao.GetVideoByVideoId(videoId)

	//包括软删除一起查找是否存在这条记录
	if found := dao.Db.Unscoped().Where(&fav).First(&fav).Error == nil; found {
		if fav.DeletedAt.Valid { //如果有软删除记录，那么不用重新创建
			// If "DeletedAt.Valid" is true, it's deleted.
			// fav.DeletedAt.Valid = false
			dao.Db.Model(&fav).Unscoped().Where(&fav).Update("deleted_at", nil)
			log.Println("将软删除设置为了无效")

			// Increase the value of "favorite_count" by 1
			err1 := dao.UpdateFavoriteCountByVideoId(videoId)

			err2 := user.AddFavoriteCount(userId)

			err3 := user.AddTotalFavorite(video.Author.Id)

			if err1 == false || err2 != nil || err3 != nil {
				log.Println("点赞操作出现错误")
				return true, nil
			}

			return false, nil
		} else { //否则说明重复点赞
			log.Println("检测到重复点赞")
			return true, errors.New("repeat operation")
		}
	} else {
		//没有这条记录，正常执行点赞操作
		if err := dao.Db.Save(&fav).Error; err != nil {
			log.Println(err.Error())
			return false, err
		}
	}

	//数据库视频表点赞数 + 1
	err1 := dao.UpdateFavoriteCountByVideoId(videoId)

	err2 := user.AddFavoriteCount(userId)

	err3 := user.AddTotalFavorite(video.Author.Id)

	if err1 == false || err2 != nil || err3 != nil {
		log.Println("点赞操作出现错误")
		return true, nil
	}

	return false, nil
}

func UnFav(userId int64, videoId int64) (bool, error) {
	video, _ := dao.GetVideoByVideoId(videoId)
	log.Println(jsonUtils.MapToJson(video))
	log.Println("执行软删除操作")
	//fav := &dao.TableFavor{UserId: userId, VideoId: videoId}
	fav := &dao.TableFavor{}

	//如果软删除过了，就执行里面的处理，否则直接创建这条记录; 测试创建成功
	if found := dao.Db.Unscoped().
		Where("user_id = ? and video_id = ?", userId, videoId).
		First(&fav).Error == nil; found {

		log.Println("成功找到喜欢记录")
		if fav.DeletedAt.Valid == true {
			//有软删除记录，说明重复操作了
			return true, errors.New("repeat operation")

		} else {
			//有，但没有软删除过，正常删除
			log.Println("执行软删除")
			log.Println(fav)
			err := dao.Db.Where(&fav).Delete(&fav).Error
			log.Println("执行软删除成功")

			if err != nil {
				log.Println(err.Error())
				log.Println("软删除失败")
				return false, err
			}

			//数据库视频表点赞数 - 1
			dao.Db.Model(&dao.Video{}).
				Where("id = ?", videoId).
				Update("favorite_count", gorm.Expr("favorite_count - ?", 1))

			//log.Printf("uid=%v\n", userId)
			user.SubFavoriteCount(userId)

			//log.Printf("vid=%v\n", video.Author.Id)
			user.SubTotalFavorite(video.Author.Id)

			return false, nil
		}
	}

	//不然就是重复操作
	return true, errors.New("repeat operation")

}
