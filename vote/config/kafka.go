package config

import "github.com/IBM/sarama"

const (
	KafkaBrokers = "localhost:9092"
	KafkaTopic   = "vote"
)

var Producer sarama.SyncProducer

func init() {
	config := sarama.NewConfig()
	//等待所有副本确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5 //最大重试次数
	config.Producer.Return.Successes = true
	config.Version = sarama.V4_1_0_0

	var err error
	Producer, err = sarama.NewSyncProducer([]string{KafkaBrokers}, config)
	if err != nil {
		panic("fail to create producer" + err.Error())
	}
}
