package service

import (
	"encoding/json"
	"strings"
	"toupiao/config"
	"toupiao/logger"
	"toupiao/model"
	"toupiao/utils"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func Start() {
	configs := sarama.NewConfig()
	configs.Consumer.Return.Errors = true
	configs.Version = sarama.V4_1_0_0

	brokers := strings.Split(config.KafkaBrokers, ",")
	consumer, err := sarama.NewConsumer(brokers, configs)
	if err != nil {
		logger.Fatal(gin.H{"fail to init consumer": err.Error()})
	}
	defer consumer.Close()

	//订阅
	list, err := consumer.Partitions(config.KafkaTopic)
	if err != nil {
		logger.Fatal(gin.H{"fail to get list of topics": err.Error()})
	}

	//爽消费爽
	for _, part := range list {
		pc, err := consumer.ConsumePartition(config.KafkaTopic, part, sarama.OffsetNewest)
		if err != nil {
			logger.Error(gin.H{"fail to consume": err.Error()})
			continue
		}
		defer pc.AsyncClose()

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				var vote utils.Vote
				err := json.Unmarshal(msg.Value, &vote)
				if err != nil {
					logger.Error(gin.H{"fail to unmarshal msg": err.Error()})
					continue
				}
				if err := model.Save(vote.ID, vote.UserID, vote.IP, vote.VoteTime); err != nil {
					logger.Error(gin.H{"fail to save vote": err.Error()})
				}
			}
		}(pc)
	}
	select {}
}
