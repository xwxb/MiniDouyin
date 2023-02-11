package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
)

type UserListResponse struct {
	Response
	UserList []dao.TableUser `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		err = nil
	}
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		err = nil
	}
	fmt.Printf("to_user_id = %d, userId = %d\n", toUserId, userId)

	ok := (len(token) > 0)

	if ok {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Fatal: invalid user id"})
	}
	userList, err := dao.GetFollowListByFollowerId(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList
func FollowerList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Fatal: invalid user id"})
	}
	userList, err := dao.GetFollowerListByFollowId(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
	fmt.Printf("c: %v\n", c)
}

// beta: same with follower list
func FriendList(c *gin.Context) {
	FollowerList(c)
}
