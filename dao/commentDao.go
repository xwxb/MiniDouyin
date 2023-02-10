package dao

import (
	"log"
)

// TableComment 评论表
type TableComment struct {
	// 评论id
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id,omitempty"`
	// (对应)视频id
	VideoId int64 `gorm:"column:video_id" json:"video_id,omitempty"`
	// (发布评论的)用户ID
	UserId int64 `gorm:"column:user_id" json:"user_id,omitempty"`
	// 评论内容
	Content string `gorm:"column:content" json:"content,omitempty"`
	// 评论发布日期, 格式 mm-dd
	CreateDate string `gorm:"column:create_date" json:"create_date,omitempty"`
	// 删除标识,true表示被删除
	Delete bool `gorm:"column:delete" json:"delete,omitempty"`
}

func (comment TableComment) TableName() string {
	return "comment"
}

// AddComment 上传评论
func AddComment(comment *TableComment) bool {
	if err := Db.Create(&comment).Error; err != nil {
		log.Println("[上传评论] 产生错误：", err)
		return false
	}
	return true
}

// DeleteComment 删除评论，参数为评论id
func DeleteComment(commentId int64) bool {
	err := Db.Model(&TableComment{Id: commentId}).Update("delete", true).Error
	if err != nil {
		log.Println("[删除评论] 产生错误：", err)
		return false
	}
	return true
}

// IsCommentUser 判断操作用户是否是发表评论的用户 (如果该 评论不存在 或 对应用户错误 会返回false)
func IsCommentUser(commentId, userId int64) bool {
	comment := TableComment{}
	// 查询评论对应发布用户
	res := Db.Select("user_id").Where(&TableComment{Id: commentId}).First(&comment)
	if res.Error != nil {
		// 没查找到记录好像也会报错，主要注意是否有其他错误
		log.Println("[判断操作用户] 产生错误：", res.Error, " (没查找到记录好像也会报错，主要注意是否有其他错误)")
		return false
	}
	// 判断用户是否正确
	return res.RowsAffected == 1 && userId == comment.UserId
}

// GetCommentList 查看视频的所有评论，按发布时间倒序 (不含删除内容)
func GetCommentList(videoId int64) (commentList []TableComment, err error) {
	var comments []TableComment
	res := Db.Model(&TableComment{}).Where(map[string]interface{}{"video_id": videoId, "delete": false}).Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}
