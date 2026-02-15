package service

import (
	"errors"

	"github.com/RBAC/internal/repository"
)

var getRoleByUserId = repository.GetRoleByUserId
var setUserRole = repository.SetUserRole
var checkUserExists = repository.UserExistsByID
var checkRoleExists = repository.RoleExistsByID

func SetUserRole(userId, roleId int64) error {
	if userId <= 0 {
		return errors.New("请选择一个用户")
	}
	if roleId <= 0 {
		return errors.New("请选择一个角色")
	}

	userExists, err := checkUserExists(userId)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("用户不存在")
	}

	roleExists, err := checkRoleExists(roleId)
	if err != nil {
		return err
	}
	if !roleExists {
		return errors.New("角色不存在")
	}

	roleIds, err := getRoleByUserId(userId)
	if err != nil {
		return err
	}

	if len(roleIds) > 0 {
		for _, assignedRoleID := range roleIds {
			if assignedRoleID == roleId {
				return errors.New("用户已拥有该角色")
			}
		}
		return errors.New("用户已分配其他角色，请先解绑")
	}

	if err := setUserRole(userId, roleId); err != nil {
		return err
	}

	return nil
}
