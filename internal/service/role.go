package service

import (
	"errors"

	"github.com/RBAC/internal/repository"
)

var getRoleByUserId = repository.GetRoleByUserId
var setUserRole = repository.SetUserRole

func SetUserRole(userId, roleId int64) error {
	if roleId == 0 {
		return errors.New("选择一个角色")
	}
	roleIds, err := getRoleByUserId(userId)
	if err != nil {
		return err
	}
	if len(roleIds) > 0 {
		return errors.New("用户已分配角色，不能重复分配")
	}

	if err := setUserRole(userId, roleId); err != nil {
		return err
	}

	return nil
}
