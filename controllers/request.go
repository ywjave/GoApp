package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

const CtxUidKey = "UserID"

func GetCurrentUserID(c *gin.Context) (UserID int64, err error) {
	uid, ok := c.Get(CtxUidKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	UserID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	fmt.Println("检测uid", uid)
	return
}

func getpageinfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var (
		page int64
		size int64
		err  error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
