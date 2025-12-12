package handler

import (
	"net/http"

	"github.com/RBAC/internal/service"
	tools "github.com/RBAC/pkg/ecode"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Account  string `json:"account", binding:"required", form:"account"`
	Password string `json:"password", binding:"required", form:"password"`
}

func DoLogin(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, tools.ParamErr)
		return
	}

	token, err := service.Login(ctx, req.Account, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, tools.UserErr)
		return
	}

	ctx.JSON(http.StatusOK, tools.ECode{
		Data: gin.H{
			"token": token,
		},
	})
}
