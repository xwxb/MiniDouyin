package module

import (
	"errors"
	"fmt"
	"github.com/xwxb/MiniDouyin/config"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

/* 因为在原来表上加 jwt.StandardClaims 好像会引发一些bug？
所以就新建了一个用来认证的类。
*/
type userStdClaims struct {
	jwt.StandardClaims
	*dao.TableUser
}

// 根据 用户（user） 信息创建 token
func JwtGenerateToken(user *dao.TableUser, duration time.Duration) string {
	fmt.Printf("JWTuser = %v\n", user)
	expireTime := time.Now().Add(duration)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", user.Id),
		Issuer:    "github.com/xwxb/MiniDouyin",
	}

	userClaims := userStdClaims{
		StandardClaims: stdClaims,
		TableUser:      user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	if tokenString, err := token.SignedString([]byte(config.SecretKey)); err == nil {
		tokenString = "Bearer " + tokenString
		println("TOKENSTRING = ", tokenString)
		println("generate token success!\n")
		return tokenString
	} else {
		fmt.Printf("generate token fail : %v\n", err)
		return "failed to generate token"
	}
}

// 将用户信息从 token 中解析出来
func JwtParseUser(tokenString string) (*dao.TableUser, error) {
	if tokenString == "" {
		return nil, errors.New("no token is found")
	}

	claims := userStdClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return claims.TableUser, err
}

// 注册密码加密
func Encoder(password string) string {
	encoder, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	fmt.Printf("Encoder Result: %v\n", encoder)
	return string(encoder)
}
