package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/module"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	// (通过token验证解析后得到的)user(的引用)
	userP := c.MustGet("authUserObj").(*dao.TableUser)
	// 视频id (注意获取到的是string，要转成int64)
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	// 1-发布评论，2-删除评论
	actionType := c.Query("action_type")

	if actionType == "1" { // 发布评论
		// 用户填写的评论内容，在action_type=1的时候使用
		commentText := c.Query("comment_text")
		if len(commentText) > 256 {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  "评论字数过多",
				"comment":     nil,
			})
			return
		}
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
				"comment":     nil,
			})
		}
	} else if actionType == "2" {
		// 要删除的评论id，在action_type=2的时候使用
		commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		dao.DeleteComment(commentId)
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	// token (可能不存在或已过期)
	token := c.Query("token")
	// 视频id
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	// 根据token获取请求方的user
	userP, err := module.JwtParseUser(strings.Fields(token)[1])
	if err != nil {
		log.Printf("%+v\n", err)
	}
	// 获取评论列表
	list, err := dao.GetCommentList(videoId)
	if err != nil {
		log.Printf("[获取评论列表] 产生异常: %+v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"status_code":  1,
			"status_msg":   "数据异常",
			"comment_list": nil,
		})
		return
	}

	// TableComment中只有user_id，没有user完整信息和是否关注，需要转换和查找
	commentList := toComment(list, userP)
	c.JSON(http.StatusOK, gin.H{
		"status_code":  0,
		"status_msg":   "成功",
		"comment_list": commentList,
	})
}

// toComment 将TableComment转换为Comment对象，他们的区别在于TableComment与数据库直接对应，不包含User完整信息
func toComment(tableCommentList []dao.TableComment, user *dao.TableUser) []Comment {
	var res []Comment
	if user == nil { // 如果没有请求方用户信息(未登录或token过期)，不获取is_follow字段
		for _, c := range tableCommentList {
			tableUser, _ := dao.GetUserByUserId(c.UserId)
			u := User{
				Id:            c.UserId,
				Name:          tableUser.UserName,
				Password:      "",
				FollowCount:   tableUser.FollowCount,
				FollowerCount: tableUser.FollowerCount,
				IsFollow:      false,
			}
			res = append(res, Comment{
				Id:         c.Id,
				User:       u,
				Content:    c.Content,
				CreateDate: c.CreateDate,
			})
		}

	} else { // 请求方已登录(token解析出user)，要获取is_follow字段
		//先获取请求方用户id，再遍历进行转换和查找
		postUserId := user.Id
		for _, c := range tableCommentList {
			tableUser, _ := dao.GetUserByUserId(c.UserId)
			isFollowed, _ := dao.IsFollowed(postUserId, c.UserId)
			u := User{
				Id:            c.UserId,
				Name:          tableUser.UserName,
				Password:      "",
				FollowCount:   tableUser.FollowCount,
				FollowerCount: tableUser.FollowerCount,
				IsFollow:      isFollowed,
			}
			res = append(res, Comment{
				Id:         c.Id,
				User:       u,
				Content:    c.Content,
				CreateDate: c.CreateDate,
			})
		}
		//fmt.Printf("用户登录情况下 返回结果: %+v", res)
	}
	return res
}
