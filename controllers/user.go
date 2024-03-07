package controllers

import (
	"GoApp/dao/mysql"
	"GoApp/logic"
	"GoApp/models"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamsSignup)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithParams(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//手动校验请求参数，与业务规则相同
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
	//	zap.L().Error("SignUp with param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}
	//fmt.Println(p)
	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("Logic error", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExit) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//1.获取参数并进行校验
	p := new(models.ParamsLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithParams(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Login error", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExit) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	//3.返回响应
	//一切正常时可以返回token
	ResponseSuccess(c, gin.H{
		"user_name": user.Username,
		"user_id":   fmt.Sprintf("%d", user.UserID), //前端js id<2^53,go int64大于>2^53
		"token":     user.Token,
	})
}
