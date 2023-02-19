package user

import (
	"github.com/xwxb/MiniDouyin/dao"
	"log"
)

func AddTotalFavorite(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "total_favorited", 1); err != nil {
		log.Println("增加total_favorited失败")
		return err
	}
	return nil
}

func SubTotalFavorite(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "total_favorited", -1); err != nil {
		log.Println("减少total_favorited失败")
		return err
	}
	return nil
}

func AddFavoriteCount(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "favorite_count", 1); err != nil {
		log.Println("增加favorite_count失败")
		return err
	}
	return nil
}

func SubFavoriteCount(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "favorite_count", -1); err != nil {
		log.Println("减少favorite_count失败")
		return err
	}
	return nil
}

func AddFollowCount(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "follow_count", 1); err != nil {
		log.Println("增加follow_count失败")
		return err
	}
	return nil
}

func SubFollowCount(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "follow_count", -1); err != nil {
		log.Println("减少follow_count失败")
		return err
	}
	return nil
}

func AddFollowerCount(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "follower_count", 1); err != nil {
		log.Println("增加follower_count失败")
		return err
	}
	return nil
}

func SubFollowerCount(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "follower_count", -1); err != nil {
		log.Println("减少follower_count失败")
		return err
	}
	return nil
}

func AddWorkCount(userId int64) error {
	if err := dao.UpdateUserByNum(userId, "work_count", 1); err != nil {
		log.Println("增加work_count失败")
		return err
	}
	return nil
}
