package repository

func GetRoleByUserId(userId int64) ([]int64, error) {
	var roleIds = []int64{}
	err := Conn.Table("user_role").Where("user_id = ?", userId).Pluck("role_id", &roleIds).Error

	if err != nil {
		return nil, err
	}
	return roleIds, nil
}
