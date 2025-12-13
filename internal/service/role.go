package service

import (
	"errors"

	"github.com/RBAC/internal/repository"
)

func SetUserRole(userId, roleId int64) error {
	// Implementation for setting a user's role
	if roleId == 0 {
		return errors.New("选择一个角色")
	}
	exist, err := repository.GetRoleByUserId(userId)
	if err != nil {
		return err
	}
	if exist != nil {
		return errors.New("用户已分配角色，不能重复分配")
	}

	repository.SetUserRole(userId, roleId)

	return nil
}
