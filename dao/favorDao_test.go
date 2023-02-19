package dao

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestGetFavorList(t *testing.T) {
	var favorVideo []TableVideo
	for i := 1; i < 5; i++ {
		// t.Log.Println('-')
		t.Log()
	}
	// 测试 函数 getfavorvideoinfolistbyuserid
	userFavorInfo, err := GetFavorVideoInfoListByUserId(1)
	if err != nil {
		log.Println(err)
	}
	jsonErr := json.Unmarshal([]byte(userFavorInfo), &favorVideo)
	if jsonErr != nil {
		log.Println("favor_test.getFavorVideoInfoListByUserId is running")
		log.Println("解码失败")
	}
	fmt.Println(userFavorInfo)

	// 测试 getfavorlistbyuserid 函数
	favorList, err := GetFavorListByUserId(1)
	if err != nil {
		log.Println(err)
		log.Println("favorlist is wake")
	}
	// 把得到的 结构体数组 favorlist 中的信息打印出来观察
	for index, value := range favorList {
		log.Println("index = ", index, "videoName = ", value.Title)
	}
	for i := 1; i < 5; i++ {
		fmt.Println('-')
	}
}

