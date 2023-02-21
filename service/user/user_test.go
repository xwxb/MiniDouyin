package user

import (
	"fmt"
	"testing"
)

func TestAddFavoriteCount(t *testing.T) {
	if err := AddFavoriteCount(1); err != nil {
		fmt.Println("增加favorite_count失败")
		fmt.Println(err)
	}
	fmt.Println("增加favorite_count成功")
}

func TestSubFavoriteCount(t *testing.T) {
	if err := SubFavoriteCount(1); err != nil {
		fmt.Println("减少favorite_count失败")
		fmt.Println(err)
	}
	fmt.Println("减少favorite_count成功")
}

func TestAddFollowCount(t *testing.T) {
	if err := AddFollowCount(1); err != nil {
		fmt.Println("增加follow_count失败")
		fmt.Println(err)
	}
	fmt.Println("增加follow_count成功")
}

func TestSubFollowCount(t *testing.T) {
	if err := SubFollowCount(1); err != nil {
		fmt.Println("减少follow_count失败")
		fmt.Println(err)
	}
	fmt.Println("减少follow_count成功")
}

func TestAddFollowerCount(t *testing.T) {
	if err := AddFollowerCount(1); err != nil {
		fmt.Println("增加follower_count失败")
		fmt.Println(err)
	}
	fmt.Println("增加follower_count成功")
}

func TestSubFollowerCount(t *testing.T) {
	if err := SubFollowerCount(3); err != nil {
		fmt.Println("减少follower_count失败")
		fmt.Println(err)
	}
	fmt.Println("减少follower_count成功")
}

func TestAddWorkCount(t *testing.T) {
	if err := AddWorkCount(1); err != nil {
		fmt.Println("增加work_count失败")
		fmt.Println(err)
	}
	fmt.Println("增加work_count成功")
}

func TestAddTotalFavorite(t *testing.T) {
	if err := AddTotalFavorite(1); err != nil {
		fmt.Println("增加total_favorite失败")
		fmt.Println(err)
	}
	fmt.Println("增加total_favorite成功")
}

func TestSubTotalFavorite(t *testing.T) {
	if err := SubFollowerCount(1); err != nil {
		fmt.Println("减少total_favorite失败")
		fmt.Println(err)
	}
	fmt.Println("减少total_favorite成功")
}
