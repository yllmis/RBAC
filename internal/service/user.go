package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RBAC/internal/model"
	"github.com/RBAC/internal/repository"
	"github.com/RBAC/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func Register(ctx context.Context, name, account, password string) error {
	_ = ctx
	name = strings.TrimSpace(name)
	account = strings.TrimSpace(account)
	password = strings.TrimSpace(password)

	if account == "" || password == "" {
		return errors.New("账号或密码不能为空")
	}

	exists, err := repository.AccountExists(account)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("账号已存在")
	}

	return repository.CreateUser(name, account, password)
}

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
	if err != nil {
		return "", err
	}

	roleID, err := repository.GetRoleByUserId(user.Id)
	if err != nil {
		return "", err
	}
	perms, err := repository.GetPermByRole(roleID)
	if err != nil {
		return "", err
	}

	redisKey := fmt.Sprintf("user_perms_%d", user.Id)
	permMap := make(map[string]interface{})

	for _, p := range perms {
		if p.PermCode != "" {
			permMap[p.PermCode] = "1"
		}
	}

	if repository.Rdb != nil && len(permMap) > 0 {
		_, pipeErr := repository.Rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, redisKey, permMap)
			pipe.Expire(ctx, redisKey, utils.TokenTTL)
			return nil
		})
		if pipeErr != nil {
			fmt.Printf("Redis缓存失败,%s", pipeErr)
		}
	}

	return token, err
}

func GetUserList(ctx context.Context) ([]model.User, error) {
	return repository.GetAllUsers()
}
