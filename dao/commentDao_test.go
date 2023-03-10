package dao

import (
	"fmt"
	"log"
	"testing"
)

// 上传评论测试
func TestAddComment(t *testing.T) {
	comment := TableComment{
		VideoId:    1,
		UserId:     6,
		Content:    "这是一条评论",
		CreateDate: "02-10",
		Delete:     false,
	}
	AddComment(&comment)
	// 查找数据库查看结果
	Db.First(&comment)
	fmt.Printf("\nid = %d\n", comment.Id)
}

// 删除评论测试
func TestDeleteComment(t *testing.T) {
	DeleteComment(3)
	comment := TableComment{Id: 3}
	Db.First(&comment)
	fmt.Printf("\ncomment after  = %+v\n", comment)
}

// 判断是否为评论发布人测试
func TestIsCommentUser(t *testing.T) {
	isCommentUser := IsCommentUser(2, 6)
	fmt.Printf("是否为发布该评论的用户: %v\n", isCommentUser)
	isCommentUser = IsCommentUser(2, 1)
	fmt.Printf("是否为发布该评论的用户: %v\n", isCommentUser)
}

// 获取视频评论
func TestGetCommentList(t *testing.T) {
	commentList, err := GetCommentList(1)
	if err != nil {
		log.Printf("[获取评论异常] 异常%+v", err)
	}
	for idx, comment := range commentList {
		fmt.Printf("%d - %+v\n", idx, comment)
	}
}

func TestGetCommentNum(t *testing.T) {
	num, err := GetCommentNum(1)
	if err != nil {
		log.Printf("[获取评论数量异常] 异常%+v", err)
	}
	fmt.Printf("评论数:%d\n", num)

	num, err = GetCommentNum(2)
	if err != nil {
		log.Printf("[获取评论数量异常] 异常%+v", err)
	}
	fmt.Printf("评论数:%d\n", num)
}
