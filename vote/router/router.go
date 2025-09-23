package router

import (
	"net/http"
	"toupiao/config"
	"toupiao/controllers"
	"toupiao/logger"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default() //整个默认路由
	router.Use(logger.LoggerToFile())
	router.Use(logger.Recover)

	rds, _ := redis.NewStore(10, "tcp", config.RedisAddress, "xxh2023gkpku", "secret")
	router.Use(sessions.Sessions("mysession", rds))

	//该定义分组了
	//这样才能在url为user中访问

	user := router.Group("/user")
	{
		user.GET("/info/:id", controllers.UserController{}.GetUserInfo)
		user.POST("/list", controllers.UserController{}.GetList)
		user.POST("/add", controllers.UserController{}.AddUser)
		user.POST("/update", controllers.UserController{}.UpdateUser)
		user.POST("/delete", controllers.UserController{}.DeleteUser)
		user.POST("/listtest", controllers.UserController{}.GetUserListTest)
		user.PUT("/put", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "put")
		})
		user.DELETE("/delete", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "delete")
		})
	}

	return router
}
