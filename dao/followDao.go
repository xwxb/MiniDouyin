package dao

import (
	"log"

	"gorm.io/gorm"
)

type Follow struct {
	Id         int64 `gorm:"primary_key;AUTO_INCREMENT"`
	FollowId   int64 `gorm:"follow_id"`
	FollowerId int64 `gorm:"follower_id"`
	gorm.DeletedAt
}

func (Follow) TableName() string {
	return "follow"
}

// IsFollowed(A, B) returns if A is followed by B
func IsFollowed(followId, followerId int64) (bool, error) {
	var followList []Follow

	condi := "follow_id = ? AND follower_id = ?"
	if err := Db.Where(condi, followId, followerId).Find(&followList).Error; err != nil {
		log.Println(err.Error())
	}
	return (len(followList) != 0), nil
}

// Get all users that the given specific follows.
//
// Given id of user A, returns list of all user B satisfying "A follows B".
func GetFollowListByFollowerId(followerId int64) ([]TableUser, error) {
	var followList []TableUser

	condi := "JOIN follow ON follow_id = user.id AND follower_id = ? AND deleted_at IS NULL"
	if err := Db.Joins(condi, followerId).Find(&followList).Error; err != nil {
		log.Println(err.Error())
	}
	return followList, nil
}

// Get all users who follow the specific user.
//
// Given id of user B, returns list of all user B satisfying "A follows B".
func GetFollowerListByFollowId(followId int64) ([]TableUser, error) {
	var followerList []TableUser

	condi := "JOIN follow ON follower_id = user.id AND follow_id = ? AND deleted_at IS NULL"
	if err := Db.Joins(condi, followId).Find(&followerList).Error; err != nil {
		log.Println(err.Error())
	}
	return followerList, nil
}

// UpFollow(A, B) makes A followed by B
func UpFollow(followId, followerId int64) (bool, error) {
	value := &Follow{FollowId: followId, FollowerId: followerId}
	assign := &Follow{DeletedAt: gorm.DeletedAt{Valid: false}}
	follow := &Follow{}

	if err := Db.Unscoped().Where(&value).Assign(&assign).FirstOrCreate(&follow).Error; err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
}

// Unfollow(A, B) makes A unfollowed by B
func Unfollow(followId, followerId int64) (bool, error) {
	follow := &Follow{FollowId: followId, FollowerId: followerId}

	err := Db.Where(&follow).Delete(&follow).Error
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
}
