package handler

import (
	"net/http"

	"github.com/RBAC/internal/model"
	"github.com/RBAC/internal/service"
	tools "github.com/RBAC/pkg/ecode"
	"github.com/gin-gonic/gin"
)

func SetRole(ctx *gin.Context) {

	var userRole model.UserRole
	if err := ctx.ShouldBind(&userRole); err != nil {
		ctx.JSON(http.StatusBadRequest, tools.ParamErr)
		return
	}

	err := service.SetUserRole(userRole.UserId, userRole.RoleId)
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
