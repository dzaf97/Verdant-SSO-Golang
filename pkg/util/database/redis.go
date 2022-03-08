package database

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var RedisObject *redis.Client

func NewRedis() (error, bool) {

	host, _ := os.LookupEnv("REDIS_DB_ADDR")
	password, _ := os.LookupEnv("REDIS_DB_PASSWD")

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       4,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Println(err)
		return err, false
	} else {
		RedisObject = client
	}

	log.Println(pong)

	return nil, true
}

func GetRedisInstance() *redis.Client {
	return RedisObject
}
