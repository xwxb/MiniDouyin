package module

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type userStdClaims struct {
	jwt.StandardClaims
	dao.TableUser
}

// 根据 用户（user） 信息创建 token
func JwtGenerateToken(user dao.TableUser, duration time.Duration) string {
	user.Password = ""
	expireTime := time.Now().Add(duration)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", user.Id),
		Issuer:    "github.com/RaymondCode/simple-demo",
	}

	userClaims := userStdClaims{
		StandardClaims: stdClaims,
		TableUser:      user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, userClaims)

	if tokenString, err := token.SignedString([]byte("")); err != nil {
		tokenString = "Bearer " + tokenString
		println("generate token success!\n")
		return tokenString
	} else {
		println("generate token fail\n")
		return "failed to generate token"
	}
}

// 解析 token
func JwtParseUser(tokenString string) (dao.TableUser, error) {
	if tokenString == "" {
		return dao.TableUser{}, errors.New("no token is found")
	}

	claims := userStdClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(""), nil
	})
	if err != nil {
		return dao.TableUser{}, err
	}
	return claims.TableUser, err
}

// 注册密码加密
func Encoder(password string) string {
	method := hmac.New(sha256.New, []byte(password))
	encoder := hex.EncodeToString(method.Sum(nil))
	fmt.Println("Encoder Result: " + encoder)
	return encoder
}
