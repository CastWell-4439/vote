package config

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	RedisAddress = "localhost:6379"
)

var RedisPool *redis.Pool

func init() {
	RedisPool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   20,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", RedisAddress)
			if err != nil {
				return nil, err
			} //goland真有一键if err！=nil，鱼鱼蒸好了
			return conn, nil
		},
	}
}
