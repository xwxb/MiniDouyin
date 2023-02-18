package dao

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type TableUser struct {
	Id              int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id,omitempty"`
	UserName        string `gorm:"column:user_name" json:"name,omitempty"`
	Password        string `gorm:"column:password" json:"-"` // 转化为json格式的时候自动忽略password字段
	FollowCount     int64  `gorm:"column:follow_count" json:"follow_count"`
	FollowerCount   int64  `gorm:"column:follower_count" json:"follower_count,omitempty"`
	IsFollow        bool   `gorm:"-" json:"is_follow"`
	Avatar          string `gorm:"avatar" json:"avatar"`
	BackGroundImage string `gorm:"background_image" json:"background_image"`
	Signature       string `gorm:"signature" json:"signature"`
	TotalFavorited  int    `gorm:"total_favorited" json:"total_favorited"`
	WorkCount       int    `gorm:"work_count" json:"work_count"`
	FavoriteCount   int    `gorm:"favorite_count" json:"favorite_count"`
}

type User struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	Name          string `json:"name,omitempty" gorm:"column:user_name"`
	Password      string `json:"-" gorm:"column:password"`
	FollowCount   int64  `json:"follow_count,omitempty" gorm:"column:follow_count"`
	FollowerCount int64  `json:"follower_count,omitempty" gorm:"column:follower_count"`
	IsFollow      bool   `json:"is_follow,omitempty" gorm:"-"`
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

// UpdateUserByNum
//
//	@Description: 更新 User 的数字字段（包括视频数量，赞数量，获赞数量等等）
//	@param userId
//	@param column 字段名字
//	@param oper 操作类型：以点赞为例为 1， 取消赞为 -1
//	@return bool
func UpdateUserByNum(userId int64, column string, oper int64) error {
	// 如果是增加字段，直接执行
	if oper == 1 {
		if err := Db.Model(&TableUser{Id: userId}).Update(column, gorm.Expr(column+"+ ?", oper)).Error; err != nil {
			log.Printf("更新%v的%v失败", userId, column)
			return err
		}
		return nil
	}

	// 如果是减少字段，那么需要开启事务，防止减少后数量小于 0。
	// 开启事务后，先减少字段的数，然后判断该字段是否小于 0。
	tx := Db.Begin()
	if err := tx.Model(&TableUser{Id: userId}).Update(column, gorm.Expr(column+"+ ?", oper)).Error; err != nil {
		log.Printf("更新%v的%v失败", userId, column)
		tx.Rollback()
		return err
	}

	var num int64

	if tx.Table("user").Where("id = ? ", userId).Select(column).Scan(&num); num < 0 {
		tx.Rollback()
		return errors.New("不能少于0")
	}

	tx.Commit()
	return nil
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
