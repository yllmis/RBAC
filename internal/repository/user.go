package repository

import (
	"fmt"

	"github.com/RBAC/internal/model"
)

func GetUserByAccount(account string) (*model.User, error) {
	var user model.User

	err := Conn.Raw("select * from user where account = ? limit 1", account).Scan(&user).Error
	if err != nil {
		return nil, fmt.Errorf("GetUserByAccount error: %v", err)
	}

	return &user, nil
}
