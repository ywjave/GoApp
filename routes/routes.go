package routes

import (
	"GoApp/controllers"
	"GoApp/logger"
	"GoApp/middlewares"
	"net/http"
	_ "strings"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	if viper.GetString("app.mode") == "realse" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//r.LoadHTMLGlob("./templates/index.html")
	//r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	v1 := r.Group("/api/v1")
	v1.POST("/signup", controllers.SignUpHandler)
	v1.POST("/login", controllers.LoginHandler)
	//时间或分数进行排序
	v1.GET("/posts2/", controllers.GetPostListHandler2) //浏览无需登录
	v1.Use(middlewares.JWTAuthMiddleware())             //使用中间件
	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)
		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		v1.GET("/posts", controllers.GetPostListHandler)
		v1.POST("/vote", controllers.PostVoteHandler)
	}
	//pprof.Register(r) //注册pprof相关路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404 ",
		})
	})
	return r
}

// JWTAuthMiddleware 基于JWT的认证中间件
