package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RBAC/internal/model"
	"github.com/RBAC/internal/repository"
	"github.com/RBAC/pkg/utils"
)

func Login(ctx context.Context, account, password string) (string, error) {

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
	roleId, err := repository.GetRoleByUserId(user.Id)
	perms, err := repository.GetPermByRole(roleId)
	if err != nil {
		return "", err
	}

	rediskey := fmt.Sprintf("user_perms_%d", user.Id)
	permMap := make(map[string]interface{})

	for _, p := range perms {
		// 确保 PermCode 不为空
		if p.PermCode != "" {
			// Field: "user:list", Value: "1"
			permMap[p.PermCode] = "1"
		}

		if p.Method != "" && p.ApiPath != "" {
			key := fmt.Sprintf("%s:%s", p.Method, p.ApiPath)
			permMap[key] = "1"
		}
	}

	if len(permMap) > 0 {
		err = repository.Rdb.HSet(ctx, rediskey, permMap).Err()
		if err != nil {
			fmt.Printf("Redis缓存失败,%s", err)
		}

		repository.Rdb.Expire(ctx, rediskey, 2*time.Hour)
	}

	return token, err

}

func GetUserList(ctx context.Context) ([]model.User, error) {
	return repository.GetAllUsers()
}
