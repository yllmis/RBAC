package repository

import (
	"fmt"
	"testing"
)

func TestGetUserAccounts(t *testing.T) {
	NewMysql()

	user, _ := GetUserByAccount("admin")
	fmt.Printf("user: %+v", user)

}

func TestGetRoleById(t *testing.T) {
	NewMysql()

	roles, _ := GetRoleByUserId(1)
	for _, role := range roles {
		fmt.Printf("role: %+v\n", role)
	}
}

func TestSetUserRole(t *testing.T) {
	NewMysql()

	err := SetUserRole(2, 2)
	if err != nil {
		t.Errorf("SetUserRole failed: %v", err)
	}
}
