package controllers

import (
	"strconv"
	"toupiao/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
}

func (u UserController) Register(ctx *gin.Context) {
	username := ctx.DefaultPostForm("username", "")
	password := ctx.DefaultPostForm("password", "")
	surepassword := ctx.DefaultPostForm("surepassword", "")

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ReturnError(ctx, 400, gin.H{"error": "can't use this password"})
		return
	}
	if username == "" || password == "" || surepassword == "" {
		ReturnError(ctx, 400, gin.H{"error": "mistake information"})
		return
	}

	if password != surepassword {
		ReturnError(ctx, 400, gin.H{"error": "make sure you input same password"})
		return
	}

	user, err := model.GetUserData(username)

	if user.Id != 0 {
		ReturnError(ctx, 400, gin.H{"error": "the user already exist"})
		return
	}

	ifm, err := model.AddUser(username, string(hashpassword))

	if err != nil {
		ReturnError(ctx, 400, gin.H{"error": "fail to create user"})
		return
	}

	ReturnSuccess(ctx, 200, gin.H{"ok": "OK"}, ifm, 1)
}

func (u UserController) Login(ctx *gin.Context) {
	username := ctx.DefaultPostForm("username", "")
	password := ctx.DefaultPostForm("password", "")

	if username == "" || password == "" {
		ReturnError(ctx, 4001, gin.H{"error": "input username or password"})
		return
	}

	ifm, _ := model.GetUserData(username)

	err := bcrypt.CompareHashAndPassword([]byte(ifm.Password), []byte(password))
	if err != nil {
		ReturnError(ctx, 4004, gin.H{"error": "username or password wrong"})
		return
	}

	session := sessions.Default(ctx) //把用户信息存一下
	session.Set("login:"+strconv.Itoa(ifm.Id), ifm.Id)
	session.Save()
	ReturnSuccess(ctx, 200, gin.H{"ok": "username or password wrong"}, ifm, 1)
}
