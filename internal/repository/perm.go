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
