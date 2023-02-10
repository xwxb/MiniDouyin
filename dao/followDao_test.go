package dao

import (
	"fmt"
	"testing"
)

func TestJudgeFavorByUserId(t *testing.T) {
	Init()
	isFavorite := JudgeFavorByUserId(3, 1)
	fmt.Println(isFavorite)
}
