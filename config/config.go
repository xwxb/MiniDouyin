package config

import "time"

var SecretKey = "github.com/RaymondCode/simple-demo" // jwt 加密字符串
var Duration = time.Hour * 24 * 365                  // jwt 认证时间

var LoginHead = "root:114514@tcp(47.94.10.223:3306)/mdy"
