package controllers

import (
	"strconv"
	"toupiao/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (u UserController) Register(ctx *gin.Context) {
	username := ctx.DefaultPostForm("username", "")
	password := ctx.DefaultPostForm("password", "")
	surepassword := ctx.DefaultPostForm("surepassword", "")

	if username == "" || password == "" || surepassword == "" {
		ctx.JSON(400, gin.H{"error": "mistake information"})
		return
	}

	if password != surepassword {
		ctx.JSON(400, gin.H{"error": "make sure you input same password"})
		return
	}

	user, err := model.GetUserData(username)

	if user.Id != 0 {
		ctx.JSON(400, gin.H{"error": "the user already exist"})
		return
	}

	_, err = model.AddUser(username)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "fail to create user"})
		return
	}

	ctx.JSON(200, gin.H{"succes": "OK"})
}

func (u UserController) Login(ctx *gin.Context) {
	username := ctx.DefaultPostForm("username", "")
	password := ctx.DefaultPostForm("password", "")

	if username == "" || password == "" {
		ctx.JSON(4001, "input username or password")
		return
	}

	ifm, _ := model.GetUserData(username)

	if password != ifm.Password {
		ctx.JSON(4004, "username or password wrong")
	}

	session := sessions.Default(ctx) //把用户信息存一下
	session.Set("login:"+strconv.Itoa(ifm.Id), ifm.Id)
	session.Save()
	ctx.JSON(200, "perfect")
}
