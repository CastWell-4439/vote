package model

import (
	"toupiao/config"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "user"
}

// 查询
func GetUserData(id string) (User, error) {
	var user User
	err := config.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

func GetUserListTest() ([]User, error) {
	var users []User
	err := config.DB.Where("id < ?", 3).Find(&users).Error
	return users, err
}

func AddUser(username string, password string) (int, error) {
	user := User{
		Username: username,
		Password: password,
	}
	err := config.DB.Create(&user).Error
	return user.Id, err
}

func UpdateUser(id int, username string) {
	config.DB.Model(&User{}).Where("id = ?", id).Update("username", username)
}

func DeleteUser(id int) error {
	err := config.DB.Delete(&User{}, id).Error
	return err
}
