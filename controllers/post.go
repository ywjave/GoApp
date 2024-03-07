package controllers

import (
	"GoApp/logic"
	"GoApp/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	//参数校验
	post := new(models.Post)
	if err := c.ShouldBindJSON(post); err != nil {
		zap.L().Debug("c.ShouldBindJSON(post) error", zap.Any("err", err))
		zap.L().Error("create post with invaild params")

		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c取到当前发帖用户的ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	post.AuthorID = userID
	//创建帖子
	if err := logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost(post) failed ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	//解析参数
	pidstr := c.Param("id")
	pid, err := strconv.ParseInt(pidstr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//查询数据库
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := getpageinfo(c)
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func GetPostListHandler2(c *gin.Context) {
	//获取分页参数
	//page, size := getpageinfo(c)
	p := &models.ParamsPostList{
		-1,
		1,
		10,
		models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("getpostlisthandler2 with invalid parms ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostListNew(p)
	if err != nil {
		zap.L().Error("logic.GetPostList()2 failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

//func GetCommunityPostList(c *gin.Context) {
//	p := &models.ParamsCommunityPostList{
//		&models.ParamsPostList{
//			-1,
//			1,
//			10,
//			models.OrderTime,
//		},
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("getCommunitypostlisthandler2 with invalid parms ", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetCommunityPostList()2 failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	ResponseSuccess(c, data)
//}
