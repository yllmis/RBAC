package service

import (
	"errors"
	"testing"
)

func mockRoleDeps(t *testing.T) {
	t.Helper()

	originalGet := getRoleByUserId
	originalSet := setUserRole
	originalUserExists := checkUserExists
	originalRoleExists := checkRoleExists

	getRoleByUserId = func(int64) ([]int64, error) { return []int64{}, nil }
	setUserRole = func(int64, int64) error { return nil }
	checkUserExists = func(int64) (bool, error) { return true, nil }
	checkRoleExists = func(int64) (bool, error) { return true, nil }

	t.Cleanup(func() {
		getRoleByUserId = originalGet
		setUserRole = originalSet
		checkUserExists = originalUserExists
		checkRoleExists = originalRoleExists
	})
}

func TestSetUserRole_EmptyUserID(t *testing.T) {
	mockRoleDeps(t)
	err := SetUserRole(0, 2)
	if err == nil || err.Error() != "请选择一个用户" {
		t.Fatalf("expected 请选择一个用户, got %v", err)
	}
}

func TestSetUserRole_EmptyRoleID(t *testing.T) {
	mockRoleDeps(t)
	err := SetUserRole(1, 0)
	if err == nil || err.Error() != "请选择一个角色" {
		t.Fatalf("expected 请选择一个角色, got %v", err)
	}
}

func TestSetUserRole_UserNotFound(t *testing.T) {
	mockRoleDeps(t)
	checkUserExists = func(int64) (bool, error) { return false, nil }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "用户不存在" {
		t.Fatalf("expected 用户不存在, got %v", err)
	}
}

func TestSetUserRole_RoleNotFound(t *testing.T) {
	mockRoleDeps(t)
	checkRoleExists = func(int64) (bool, error) { return false, nil }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "角色不存在" {
		t.Fatalf("expected 角色不存在, got %v", err)
	}
}

func TestSetUserRole_AlreadyHasSameRole(t *testing.T) {
	mockRoleDeps(t)
	getRoleByUserId = func(int64) ([]int64, error) { return []int64{2}, nil }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "用户已拥有该角色" {
		t.Fatalf("expected 用户已拥有该角色, got %v", err)
	}
}

func TestSetUserRole_AlreadyAssignedOtherRole(t *testing.T) {
	mockRoleDeps(t)
	getRoleByUserId = func(int64) ([]int64, error) { return []int64{3}, nil }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "用户已分配其他角色，请先解绑" {
		t.Fatalf("expected 用户已分配其他角色，请先解绑, got %v", err)
	}
}

func TestSetUserRole_QueryRoleError(t *testing.T) {
	mockRoleDeps(t)
	getRoleByUserId = func(int64) ([]int64, error) { return nil, errors.New("query failed") }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "query failed" {
		t.Fatalf("expected query failed, got %v", err)
	}
}

func TestSetUserRole_SetSuccess(t *testing.T) {
	mockRoleDeps(t)
	called := false
	setUserRole = func(int64, int64) error {
		called = true
		return nil
	}

	if err := SetUserRole(1, 2); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !called {
		t.Fatalf("expected setUserRole to be called")
	}
}

func TestSetUserRole_SetFail(t *testing.T) {
	mockRoleDeps(t)
	setUserRole = func(int64, int64) error { return errors.New("insert failed") }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "insert failed" {
		t.Fatalf("expected insert failed, got %v", err)
	}
}
