package dao

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGetPublishVideoInfoListByUserId(t *testing.T) {
	Init()
	var publicVideo []TableVideo
	userPublicInfo, err := GetPublishVideoInfoListByUserId(1)
	if err != nil {
		fmt.Println(err)
	}
	jsonErr := json.Unmarshal([]byte(userPublicInfo), &publicVideo)
	if jsonErr != nil {
		fmt.Println("解码失败")
	}
	fmt.Println(userPublicInfo)
}

func TestGetVideoByCreatedTime(t *testing.T) {
	Init()
	var publicVideo []TableVideo
	publicVideoInfo, err := GetVideoByCreatedTime(time.Now())
	if err != nil {
		fmt.Println(err)
	}
	jsonErr := json.Unmarshal([]byte(publicVideoInfo), &publicVideo)
	if jsonErr != nil {
		fmt.Println("解码失败")
	}
	fmt.Println(publicVideo)
}
