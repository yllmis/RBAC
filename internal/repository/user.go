package repository

import (
	"errors"
	"fmt"

	"github.com/RBAC/internal/model"
	"gorm.io/gorm"
)

func GetUserByAccount(account string) (*model.User, error) {
	var user model.User

	err := Conn.Where("account = ?", account).Limit(1).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByAccount error: %v", err)
	}

	return &user, nil
}

func AccountExists(account string) (bool, error) {
	var cnt int64
	if err := Conn.Table("user").Where("account = ?", account).Count(&cnt).Error; err != nil {
		return false, fmt.Errorf("AccountExists error: %v", err)
	}
	return cnt > 0, nil
}

func CreateUser(name, account, password string) error {
	err := Conn.Create(&model.User{
		Name:     name,
		Account:  account,
		Password: password,
	}).Error
	if err != nil {
		return fmt.Errorf("CreateUser error: %v", err)
	}
	return nil
}

func GetAllUsers() ([]model.User, error) {
	var users []model.User

	err := Conn.Raw("select * from user").Scan(&users).Error

	if err != nil {
		return nil, fmt.Errorf("GetAllUsers error: %v", err)
	}

	return users, nil
}
