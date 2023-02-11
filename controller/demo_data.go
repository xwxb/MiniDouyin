package controller

import (
	"github.com/xwxb/MiniDouyin/dao"
)

// 将 dao.TableUser 类型转为 controller.User 类型 (临时使用)
func toUserTemp(tableUser dao.TableUser) User {
	return User{
		Id:            tableUser.Id,
		Name:          tableUser.UserName,
		FollowCount:   tableUser.FollowCount,
		FollowerCount: tableUser.FollowerCount,
		IsFollow:      tableUser.IsFollow,
	}
}

// 包含 2 个视频，视频 id 分别为 1,2 ，视频的用户 id 分别为 1,3
//
// where `init()` calls `dao.GetUserByUserId(userId)`
var DemoVideos []Video

func init() {
	// author1, err := dao.GetUserByUserId(1)
	// if err != nil {
	// 	panic(err)
	// }
	// author2, err := dao.GetUserByUserId(3)
	// if err != nil {
	// 	panic(err)
	// }

	// DemoVideos = []Video{
	// 	{
	// 		Id:            1,
	// 		Author:        toUserTemp(author1),
	// 		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
	// 		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
	// 		FavoriteCount: 0,
	// 		CommentCount:  0,
	// 		IsFavorite:    false,
	// 		Title:         "DemoTitle1",
	// 	},
	// 	{
	// 		Id:            2,
	// 		Author:        toUserTemp(author2),
	// 		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
	// 		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
	// 		FavoriteCount: 0,
	// 		CommentCount:  0,
	// 		IsFavorite:    false,
	// 		Title:         "DemoTitle2",
	// 	},
	// }
}

var DemoComments = []Comment{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
