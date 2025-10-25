package controllers

import (
	"os"
	"time"
	"toupiao/config"
	"toupiao/logger"
	"toupiao/utils"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

var VoteScipt *redis.Script

type VoteController struct{}

func init() {
	//用你自己的
	//oops有两个vote我说怎么跑不了
	scriptPath := "/home/castwell/vote/vote/scripts/lua/counter.lua"
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		logger.Fatal(gin.H{
			"msg":   "fail to read file",
			"path":  scriptPath,
			"error": err.Error(),
		})
	}
	VoteScipt = redis.NewScript(2, string(scriptContent))
}

func (v VoteController) Vote(ctx *gin.Context) {
	itemID := ctx.DefaultPostForm("itemID", "")
	userID := ctx.DefaultPostForm("userID", "")
	if itemID == "" || userID == "" {
		ReturnError(ctx, 400, "can't be empty")
		return
	}

	itemKey := "vote:count" + itemID
	userKey := "vote:count" + userID

	conn := config.RedisPool.Get()
	defer conn.Close()

	result, err := redis.Int(VoteScipt.Do(conn, itemKey, userKey, userID))
	if err != nil {
		logger.Fatal(gin.H{
			"msg":    "fail to get vote",
			"itemID": itemID,
			"userID": userID,
			"error":  err.Error(),
		})
		ReturnError(ctx, 500, "try again")
		return
	}
	if result == 0 {
		ReturnError(ctx, 400, "can't do it twice")
	} else if result == -1 {
		ReturnError(ctx, 400, "number of canshu error")
	} else if result == 1 {
		ReturnSuccess(ctx, 0, "OK", nil, 0)
		vote := utils.Vote{
			ID:       itemID,
			UserID:   userID,
			VoteTime: time.Now(),
			IP:       ctx.ClientIP(),
		}
		go func() {
			if err := utils.SendVote(vote); err != nil {
				logger.Error(gin.H{"fail to send vote": err.Error()})
			}
		}()
	} else {
		ReturnError(ctx, 500, "unknown error")
	}
}
