package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xwxb/MiniDouyin/config"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/module"
	"golang.org/x/crypto/bcrypt"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      false,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user, _ := dao.GetUserByUsername(username)

	if username == user.UserName {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := dao.TableUser{
			UserName: username,
			Password: module.Encoder(password),
		}
		if dao.InsertUser(&newUser) == true {
			token := module.JwtGenerateToken(&newUser, config.Duration)
			log.Println("注册返回的id: ", newUser.Id)
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   newUser.Id,
				Token:    token,
			})
		} else {
			println("Insert new User Fail")
		}
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	fmt.Printf("username=%v, password=%v", username, password)

	user, exist := dao.GetUserByUsername(username)
	fmt.Printf("user= %v", user)

	if exist != nil {
		c.JSON(http.StatusOK, gin.H{
			"StatusCode_": "1",
			"error_msg":   "User Does not Exist",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	log.Printf("err=%v\n, userPassword=%v\n, pass=%v\n", err, []byte(user.Password), []byte(password))

	if err == nil {
		//fmt.Printf("JWTLOGIN:%v\n", module.JwtGenerateToken(&user, config.Duration))
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    module.JwtGenerateToken(&user, time.Hour*24*365),
		})
		log.Println("login success!!!")
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	id, _ := strconv.ParseInt(userId, 10, 64)
	token := c.Query("token")

	log.Printf("id = %v, token = %v", id, token)

	if tableUser, exist := dao.GetUserByUserId(id); exist != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		user := User{
			Id:            tableUser.Id,
			Name:          tableUser.UserName,
			FollowCount:   tableUser.FollowCount,
			FollowerCount: tableUser.FollowerCount,
			IsFollow:      tableUser.IsFollow,
		}
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}
}
