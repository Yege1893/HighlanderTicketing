package service

import (
	"fmt"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
)

func Register(accessToken string) error {
	user, err := GetUserInfo(accessToken)
	if err != nil {
		return err
	}
	err1 := CreateUser(&user)
	if err1 != nil {
		return err1
	}
	var userArray []model.User
	userArray, _ = GetAllUsers()
	fmt.Println(userArray)
	return nil
}
