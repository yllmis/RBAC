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

func CheckUserPerm(userId int64, requiredPerm string) (bool, error) {
	roleIds, err := GetRoleByUserId(userId)
	if err != nil {
		return false, err
	}
	if len(roleIds) == 0 {
		return false, nil
	}

	var cnt int64
	if err := Conn.Table("permission p").
		Joins("JOIN role_perm rp ON p.id = rp.perm_id").
		Where("rp.role_id IN (?) AND p.perm_code = ?", roleIds, requiredPerm).
		Count(&cnt).Error; err != nil {
		return false, err
	}

	return cnt > 0, nil
}
