package dao

import (
	"fmt"
	"testing"
)

func TestGetUserList(t *testing.T) {
	Init()
	list, err := GetUserList()
	fmt.Printf("userlist = %v", list)
	fmt.Printf("%v", err)
}

func TestGetUserByUsername(t *testing.T) {
	Init()
	user, err := GetUserByUsername("a")
	fmt.Printf("user = %v", user)
	fmt.Printf("%v", err)
}

func TestInsertUser(t *testing.T) {
	Init()
	newUser := &TableUser{
		UserName: "c",
		Password: "333",
	}

	res := InsertUser(newUser)
	fmt.Printf("res = %v", res)
}

func TestRemoveUserByUsername(t *testing.T) {
	Init()
	res := RemoveUserByUsername("b")
	fmt.Printf("res = %v", res)
}
