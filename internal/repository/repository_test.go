package repository

import (
	"fmt"
	"os"
	"testing"
)

func requireIntegrationDB(t *testing.T) {
	t.Helper()
	if os.Getenv("RBAC_RUN_INTEGRATION") != "1" {
		t.Skip("skip integration test: set RBAC_RUN_INTEGRATION=1 to run")
	}
}

func TestGetUserAccounts(t *testing.T) {
	requireIntegrationDB(t)
	NewMysql()

	user, _ := GetUserByAccount("admin")
	fmt.Printf("user: %+v", user)
}

func TestGetRoleById(t *testing.T) {
	requireIntegrationDB(t)
	NewMysql()

	roles, _ := GetRoleByUserId(1)
	for _, role := range roles {
		fmt.Printf("role: %+v\n", role)
	}
}

func TestSetUserRole(t *testing.T) {
	requireIntegrationDB(t)
	NewMysql()

	err := SetUserRole(2, 2)
	if err != nil {
		t.Errorf("SetUserRole failed: %v", err)
	}
}
