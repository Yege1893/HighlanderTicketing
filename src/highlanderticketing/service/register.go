package service

func Register(accessToken string) error {
	user, err := GetUserInfo(accessToken)
	if err != nil {
		return err
	}
	err1 := CreateUser(&user)
	if err1 != nil {
		return err1
	}
	return nil
}
