package handler

import (
	"net/http"

	"github.com/RBAC/internal/service"
	tools "github.com/RBAC/pkg/ecode"
	"github.com/gin-gonic/gin"
)

type SetRoleRequest struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required,gt=0"`
	RoleID int64 `json:"role_id" form:"role_id" binding:"required,gt=0"`
}

func SetRole(ctx *gin.Context) {
	var req SetRoleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, tools.ParamErr)
		return
	}

	err := service.SetUserRole(req.UserID, req.RoleID)
	if err != nil {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, tools.ECode{
		Code:    0,
		Message: "角色分配成功",
		Data:    nil,
	})
}

func GetUserList(ctx *gin.Context) {
	users, err := service.GetUserList(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: "获取所有用户列表失败",
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, tools.ECode{
		Code:    0,
		Message: "获取所有用户列表成功",
		Data:    users,
	})
}
