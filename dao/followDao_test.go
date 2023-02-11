package dao

import (
	"fmt"
	"testing"
)

func handle(err error) {
	if err != nil {
		fmt.Printf("err = %v\n", err.Error())
	}
}

func TestIsFollowed(t *testing.T) {
	follow, err := IsFollowed(1, 3)
	fmt.Printf("follow_rela = %v\n", follow)
	handle(err)
}

func printUserList(list []TableUser) {
	fmt.Printf("[info] list size = %d\n", len(list))
	for idx, val := range list {
		fmt.Printf("[%d-th] %v\n", idx, val)
	}
}

func TestGetFollowListByFollowerId(t *testing.T) {
	list, err := GetFollowListByFollowerId(6)
	if err != nil {
		handle(err)
	} else {
		printUserList(list)
	}
}

func TestGetFollowerListByFollowId(t *testing.T) {
	list, err := GetFollowerListByFollowId(1)
	if err != nil {
		handle(err)
	} else {
		printUserList(list)
	}
}

func testUpFollow(a, b int64) {
	ok, err := UpFollow(a, b)
	fmt.Printf("a = %d, b = %d, done = %v\n", a, b, ok)
	handle(err)
}

func TestUpFollow(t *testing.T) {
	for b := int64(3); b <= int64(6); b++ {
		testUpFollow(1, b)
	}
	for a := int64(3); a <= int64(5); a++ {
		testUpFollow(a, 6)
	}
}

func testUnfollow(a, b int64) {
	ok, err := Unfollow(a, b)
	fmt.Printf("a = %d, b = %d, done = %v\n", a, b, ok)
	handle(err)
}

func TestUnfollow(t *testing.T) {
	for b := int64(4); b <= int64(5); b++ {
		testUnfollow(1, b)
	}
	testUnfollow(5, 6)
}
