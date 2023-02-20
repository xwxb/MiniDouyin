package dao

import (
	"gorm.io/gorm"
	"log"

	"github.com/xwxb/MiniDouyin/utils/jsonUtils"
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

type Favor struct {
	Id      int64 `gorm:"primary_key;AUTO_INCREMENT"`
	UserId  int64 `gorm:"user_id"`
	VideoId int64 `gorm:"video_id"`
	gorm.DeletedAt
}

// 方法
func (Favor) TableName() string {
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

func JudgeFavorByUserIdMult(userId int64, vids []int64) []bool {

	var favs []TableFavor
	Db.Where("user_id = ? AND video_id IN (?)", userId, vids).Find(&favs)

	favMap := make(map[int64]bool)
	for _, like := range favs {
		favMap[like.VideoId] = true
	}

	favStats := make([]bool, len(vids))
	for i, vid := range vids {
		favStats[i] = favMap[vid]
	}

	return favStats

}

// 以favor表的形式得到favorlist
func GetFavorList() ([]Favor, error) {
	var favorList []Favor
	if err := Db.Find(&favorList).Error; err != nil {
		log.Println(err.Error())
		return favorList, err
	}
	return favorList, nil
}

// 在favor表中得到某个userid对应的所有videoid 待测试 好像没用
// get userid from favor
// get videoid by userid
// get favorlist by videoid
func GetFavorListByUserId(userId int64) ([]TableVideo, error) {
	var favorList []TableVideo
	if err := Db.
		Joins("Join favor f ON f.video_id = id").
		Where("favor.user_id = ?", userId).
		Find(&favorList).Error; err != nil {
		log.Println(err.Error())
		return favorList, err
	}
	return favorList, nil
}

// 喜欢列表实际调用的函数
func GetFavorVideoInfoListByUserId(userId int64) (string, error) {
	var favorVideo []TableVideo

	err := Db.Table("favor").
		Joins("left Join video v ON favor.video_id = v.id").
		Where("favor.user_id = ? and favor.deleted_at is null", userId).
		Scan(&favorVideo).Error

	if err != nil {
		log.Println("failed")
	}
	return jsonUtils.MapToJson(favorVideo), err
}
