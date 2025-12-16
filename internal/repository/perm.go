package repository

import (
	"github.com/RBAC/internal/model"
)

func GetPermByRole(roleId []int64) ([]model.Permission, error) {
	if len(roleId) == 0 {
		return nil, nil
	}

	var permsId []int64
	var perms []model.Permission

	err := Conn.Raw("select perm_id from role_perm where role_id in (?)", roleId).Scan(&permsId).Error

	err = Conn.Raw("select * from permission where id in (?)", permsId).Scan(&perms).Error
	if err != nil {
		return nil, err
	}
	return perms, nil
}

func CheckUserPerm(userId int64, path string, method string) (bool, error) {

	roleIds, err := GetRoleByUserId(userId)
	if err != nil {
		return false, err
	}

	perms, err := GetPermByRole(roleIds)
	if err != nil {
		return false, err
	}

	for _, perm := range perms {
		if perm.ApiPath == path && perm.Method == method {
			return true, nil
		}
	}
	return false, nil
}
