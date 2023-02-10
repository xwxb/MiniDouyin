package dao

import (
	"log"
)

type TableUser struct {
	Id            int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id,omitempty"`
	UserName      string `gorm:"column:user_name" json:"name,omitempty"`
	Password      string `gorm:"column:password" json:"-"` // 转化为json格式的时候自动忽略password字段
	FollowCount   int64  `gorm:"column:follow_count" json:"follow_count"`
	FollowerCount int64  `gorm:"column:follower_count" json:"follower_count,omitempty"`
	IsFollow      bool   `gorm:"-"`
}

func (user TableUser) TableName() string {
	return "user"
}

func GetUserList() ([]TableUser, error) {
	var usersList []TableUser
	if err := Db.Find(&usersList).Error; err != nil {
		log.Println(err.Error())
		return usersList, err
	}
	return usersList, nil
}

func GetUserByUsername(username string) (TableUser, error) {
	user1 := TableUser{}

	if err := Db.Where("user_name = ?", username).First(&user1).Error; err != nil {
		log.Println(err.Error())
		return user1, err
	}
	return user1, nil
}

func GetUserByUserId(id int64) (TableUser, error) {
	user := TableUser{}
	if err := Db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

func InsertUser(user *TableUser) bool {
	if err := Db.Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func RemoveUserByUsername(userName string) bool {
	if err := Db.Where("user_name like ?", userName).Delete(TableUser{}).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
