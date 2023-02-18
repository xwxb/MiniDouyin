package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
)

type ChatResponse struct {
	Response
	MessageList []dao.Message `json:"message_list"`
}

// MessageAction sends a message as described in Request
func MessageAction(c *gin.Context) {
	content := c.Query("content")
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		fromUserId := c.MustGet("authUserObj").(*dao.TableUser).Id
		message := dao.Message{
			FromUserId: fromUserId,
			ToUserId:   toUserId,
			Content:    content,
			CreateTime: time.Now().UnixMilli(),
		}

		if _, err := dao.SendMessage(&message); err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 0})
		}
	}
}

// MessageChat gives a message list between two users
func MessageChat(c *gin.Context) {
	fromUserId := c.MustGet("authUserObj").(*dao.TableUser).Id

	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	lastTime, err := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	msgList, err := dao.GetRecentMessageListByUserId(lastTime, toUserId, fromUserId)
	if err == nil {
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: msgList})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	}
}
