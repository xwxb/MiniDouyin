package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	/*token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id:         1,
					User:       user,
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}*/

	// (通过token验证解析后得到的)user(的引用)
	userP := c.MustGet("authUserObj").(*dao.TableUser)
	// 视频id (注意获取到的是string，要转成int64)
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	// 1-发布评论，2-删除评论
	actionType := c.Query("action_type")

	if actionType == "1" { // 发布评论
		// 用户填写的评论内容，在action_type=1的时候使用
		commentText := c.Query("comment_text")
		// 创建日期 (格式:mm-dd)
		createDate := time.Now().Format("01-02")

		comment := dao.TableComment{
			VideoId:    videoId,
			UserId:     userP.Id,
			Content:    commentText,
			CreateDate: createDate,
		}
		if dao.AddComment(&comment) { // 提交评论成功
			c.JSON(http.StatusOK, gin.H{
				"status_code": 0,
				"status_msg":  "评论提交成功",
				"comment": map[string]interface{}{
					"id":          comment.Id, // 评论id
					"user":        userP,
					"content":     commentText,
					"create_date": createDate,
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  "评论提交失败",
			})
		}
	} else if actionType == "2" {
		// 要删除的评论id，在action_type=2的时候使用
		commentId := c.Query("comment_id")
		fmt.Printf("%+v\n%s\n%s\n%+v\n", userP, videoId, actionType, commentId)
		// TODO 删除逻辑为完成
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
