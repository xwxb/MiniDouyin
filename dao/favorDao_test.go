package dao

import (
	"testing"
)

func TestGetFavorVideoInfoListByUserId(t *testing.T) {
	userid := int64(2)

	favorList, err := GetFavorVideoInfoListByUserId(userid)

	if err != nil {
		t.Error(err)
	}

	t.Log(favorList)
}

