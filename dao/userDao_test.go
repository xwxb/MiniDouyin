package dao

import (
	"fmt"
	"testing"
)

func TestGetUserList(t *testing.T) {
	list, err := GetUserList()
	fmt.Printf("userlist = %v", list)
	fmt.Printf("%v", err)
}

func TestGetUserByUsername(t *testing.T) {
	user, err := GetUserByUsername("b")
	fmt.Printf("user = %v", user)
	fmt.Printf("%v", err)
}

func TestGetUserByUserId(t *testing.T) {
	user, err := GetUserByUserId(1)
	fmt.Println(user)
	fmt.Println(err)

}

func TestInsertUser(t *testing.T) {
	newUser := &TableUser{
		UserName: "a",
		Password: "111111",
	}

	res := InsertUser(newUser)
	fmt.Printf("res = %v", res)
}

func TestRemoveUserByUsername(t *testing.T) {
	res := RemoveUserByUsername("b")
	fmt.Printf("res = %v", res)
}
