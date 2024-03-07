package controllers

import (
	"GoApp/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//社区相关

func CommunityHandler(c *gin.Context) {
	//1.查询到所有的社区（communi_id,community_name）
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList()", zap.Error(err))
		ResponseError(c, CodeServerBusy) //后端错误不对外暴露
		return
	}
	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	//1.查询到所有的社区（communi_id,community_name）
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail()", zap.Error(err))
		ResponseError(c, CodeServerBusy) //后端错误不对外暴露
		return
	}
	ResponseSuccess(c, data)
}
