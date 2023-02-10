package config

import "time"

var SecretKey = "github.com/RaymondCode/simple-demo" // jwt 加密字符串
var Duration = time.Hour * 24 * 365                  // jwt 认证时间
