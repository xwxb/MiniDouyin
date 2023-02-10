package fedd

import (
	"fmt"
	"github.com/xwxb/MiniDouyin/dao"
	"github.com/xwxb/MiniDouyin/utils/jsonUtils"
	"testing"
	"time"
)

func TestGetFeed(t *testing.T) {
	dao.Init()
	feed := GetFeed(time.Now())
	fmt.Println(jsonUtils.MapToJson(feed))
}
