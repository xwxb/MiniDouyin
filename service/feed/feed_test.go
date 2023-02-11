package feed

import (
	"fmt"
	"github.com/xwxb/MiniDouyin/utils/jsonUtils"
	"testing"
	"time"
)

func TestGetFeed(t *testing.T) {
	feed := GetFeed(time.Now())
	fmt.Println(jsonUtils.MapToJson(feed))
}

func TestGetFeedByUserId(t *testing.T) {
	feed := GetFeedByUserId(time.Now(), 1)
	fmt.Println(jsonUtils.MapToJson(feed))
}
