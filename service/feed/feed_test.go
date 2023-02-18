package feed

import (
	"fmt"
	"github.com/xwxb/MiniDouyin/utils/jsonUtils"
	"testing"
	"time"
)

func TestGetFeed(t *testing.T) {
	_, feed, _ := GetFeed(time.Now())
	fmt.Println(jsonUtils.MapToJson(feed))
}

func TestGetFeedByUserId(t *testing.T) {
	_, feed, _ := GetFeedByUserId(time.Now(), 3)
	fmt.Println(jsonUtils.MapToJson(feed))
}
