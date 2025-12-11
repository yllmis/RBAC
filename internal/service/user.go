package service

import (
	"errors"

	"github.com/RBAC/internal/repository"
	"github.com/RBAC/pkg/utils"
)

func Login(account, password string) (string, error) {

	user, err := repository.GetUserByAccount(account)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("用户不存在")
	}
	if user.Password != password {
		return "", errors.New("密码错误")
	}

	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		return "", err
	}
	return token, err

}
