package dao

import (
	"errors"

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
	err := Db.Where(condi, followId, followerId).Find(&followList).Error
	return (len(followList) != 0), err
}

func IsFollowedMult(followerId int64, followIds []int64) []bool {
	//怎么感觉这个函数签名参数列表命名反过来了。。
	var followList []Follow

	//find all 
	condi := "follower_id = ? AND follow_id IN (?)"
	Db.Where(condi, followerId, followIds).Find(&followList)


	//map userid to if-follow
	folMap := make(map[int64]bool)
	for _, fol := range followList {
		folMap[fol.FollowId] = true 
	}

	folStats := make([]bool, len(followIds))
	for i, followId := range followIds {
		folStats[i] = folStats[followId]
	}

	return folStats
}


// Given id of user A, returns list of all user B satisfying "A follows B".
func GetFollowListByFollowerId(followerId int64) ([]TableUser, error) {
	var followList []TableUser
	condi := "JOIN follow ON follow_id = user.id AND follower_id = ? AND deleted_at IS NULL"
	err := Db.Joins(condi, followerId).Find(&followList).Error
	return followList, err
}

// Given id of user B, returns list of all user B satisfying "A follows B".
func GetFollowerListByFollowId(followId int64) ([]TableUser, error) {
	var followerList []TableUser
	condi := "JOIN follow ON follower_id = user.id AND follow_id = ? AND deleted_at IS NULL"
	err := Db.Joins(condi, followId).Find(&followerList).Error
	return followerList, err
}

// UpFollow(A, B) makes A followed by B
func UpFollow(followId, followerId int64) (bool, error) {
	if followId == followerId {
		return false, errors.New("self following is illegal")
	}
	value := &Follow{FollowId: followId, FollowerId: followerId}
	assign := &Follow{DeletedAt: gorm.DeletedAt{Valid: false}}
	var follow = &Follow{}
	err := Db.Unscoped().Where(&value).Assign(&assign).FirstOrCreate(&follow).Error
	if err == nil && follow.DeletedAt.Valid {
		follow.DeletedAt.Valid = false
		err = Db.Unscoped().Save(&follow).Error
	}
	return err == nil, err
}

// Unfollow(A, B) makes A unfollowed by B
func Unfollow(followId, followerId int64) (bool, error) {
	follow := &Follow{FollowId: followId, FollowerId: followerId}
	err := Db.Where(&follow).Delete(&follow).Error
	return err == nil, err
}
