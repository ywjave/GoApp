package controllers

import (
	"GoApp/logic"
	"GoApp/models"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func PostVoteHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamsVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("validator.ValidationErrors", zap.Error(errs))
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //去除结构体标识
		zap.L().Error("removeTopStruct(errs.Translate(trans))", zap.Any("err:", errData))
		ResponseErrorWithParams(c, CodeInvalidParam, errData)
		return
	}
	userid, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	if err := logic.VoteForPost(userid, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
