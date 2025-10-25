package utils

import (
	"encoding/json"
	"time"
	"toupiao/config"
	"toupiao/logger"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type Vote struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	VoteTime time.Time `json:"vote_time"`
	IP       string    `json:"ip"`
}

func SendVote(v Vote) error {
	msgbyte, err := json.Marshal(v)
	if err != nil {
		logger.Error(gin.H{"fail to marshal msg": err.Error()})
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: config.KafkaTopic,
		Value: sarama.StringEncoder(msgbyte),
	}

	_, _, err = config.Producer.SendMessage(msg)
	if err != nil {
		logger.Error(gin.H{"fail to send msg": err.Error()})
		return err
	}
	return nil
}
