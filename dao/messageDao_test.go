package dao

import (
	"fmt"
	"github.com/xwxb/MiniDouyin/utils/jsonUtils"
	"log"
	"testing"
	"time"
)

func TestGetRecentMessageListByUserId(t *testing.T) {
	messages, err := GetRecentMessageListByUserId(1676881695231, 1, 2)
	if err != nil {
		log.Println(err)
	}
	log.Println(jsonUtils.MapToJson(messages))
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now())
}
