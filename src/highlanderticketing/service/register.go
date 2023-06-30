package service

import (
	"fmt"
)

func Register(accessToken string) error {

	user, err := GetUserInfo(accessToken)
	if err != nil {
		return err
	}
	fmt.Println(user)
	err1 := CreateUser(&user)
	if err1 != nil {
		return err1
	}
	/*var userArray []model.User
	userArray, _ = GetAllUsers()

	fmt.Println(userArray)*/
	return nil
}
