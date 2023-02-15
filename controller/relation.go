package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
)

type UserListResponse struct {
	Response
	UserList []dao.TableUser `json:"user_list"`
}

// RelationAction does Follow/Unfollow action
func RelationAction(c *gin.Context) {
	userId := c.MustGet("authUserObj").(*dao.TableUser).Id

	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	if actionType == 1 {
		_, err = dao.UpFollow(toUserId, userId)
	} else {
		_, err = dao.Unfollow(toUserId, userId)
	}

	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// FollowList gets follow list
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Fatal: invalid user id"})
	} else {
		userList, err := dao.GetFollowListByFollowerId(userId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		} else {
			c.JSON(http.StatusOK, UserListResponse{Response: Response{StatusCode: 0}, UserList: userList})
		}
	}
}

// FollowerList gets follower list
func FollowerList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Fatal: invalid user id"})
	} else {
		userList, err := dao.GetFollowerListByFollowId(userId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		} else {
			c.JSON(http.StatusOK, UserListResponse{Response: Response{StatusCode: 0}, UserList: userList})
		}
	}
}

// FriendList gets friend list (following & followed by at the same time)
func FriendList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	lastTime[userId] = int64(0)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Invalid user id"})
	} else {
		followerList, err := dao.GetFollowerListByFollowId(userId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		followList, err := dao.GetFollowListByFollowerId(userId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}

		// union two lists
		set := map[int64]struct{}{}
		for _, user := range followerList {
			set[user.Id] = struct{}{}
		}
		var userList []dao.TableUser
		for _, user := range followList {
			if _, has := set[user.Id]; has {
				userList = append(userList, user)
			}
		}

		c.JSON(http.StatusOK, UserListResponse{Response: Response{StatusCode: 0}, UserList: userList})
	}
}
