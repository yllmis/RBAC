package service

import (
	"errors"
	"testing"
)

func TestSetUserRole_EmptyRoleID(t *testing.T) {
	err := SetUserRole(1, 0)
	if err == nil || err.Error() != "选择一个角色" {
		t.Fatalf("expected 选择一个角色, got %v", err)
	}
}

func TestSetUserRole_AlreadyAssigned(t *testing.T) {
	originalGet := getRoleByUserId
	originalSet := setUserRole
	defer func() {
		getRoleByUserId = originalGet
		setUserRole = originalSet
	}()

	getRoleByUserId = func(int64) ([]int64, error) { return []int64{2}, nil }
	setUserRole = func(int64, int64) error { return nil }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "用户已分配角色，不能重复分配" {
		t.Fatalf("expected duplicate-assigned error, got %v", err)
	}
}

func TestSetUserRole_RepoError(t *testing.T) {
	originalGet := getRoleByUserId
	originalSet := setUserRole
	defer func() {
		getRoleByUserId = originalGet
		setUserRole = originalSet
	}()

	getRoleByUserId = func(int64) ([]int64, error) { return nil, errors.New("query failed") }
	setUserRole = func(int64, int64) error { return nil }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "query failed" {
		t.Fatalf("expected query failed, got %v", err)
	}
}

func TestSetUserRole_SetSuccess(t *testing.T) {
	originalGet := getRoleByUserId
	originalSet := setUserRole
	defer func() {
		getRoleByUserId = originalGet
		setUserRole = originalSet
	}()

	called := false
	getRoleByUserId = func(int64) ([]int64, error) { return []int64{}, nil }
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
	originalGet := getRoleByUserId
	originalSet := setUserRole
	defer func() {
		getRoleByUserId = originalGet
		setUserRole = originalSet
	}()

	getRoleByUserId = func(int64) ([]int64, error) { return []int64{}, nil }
	setUserRole = func(int64, int64) error { return errors.New("insert failed") }

	err := SetUserRole(1, 2)
	if err == nil || err.Error() != "insert failed" {
		t.Fatalf("expected insert failed, got %v", err)
	}
}
