package repository

import (
	"github.com/RBAC/internal/model"
)

func GetRoleByUserId(userId int64) ([]int64, error) {
	var roleIds = []int64{}
	err := Conn.Table("user_role").Where("user_id = ?", userId).Pluck("role_id", &roleIds).Error
	if err != nil {
		return nil, err
	}
	return roleIds, nil
}

func SetUserRole(userId int64, roleId int64) error {
	tx := Conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&model.UserRole{
		UserId: userId,
		RoleId: roleId,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
