package redisclient

import (
	"apple/config"
	"apple/utils"
	"github.com/go-redis/redis"
)

var RedisClient = redis.NewClient(&redis.Options{
	Addr: config.ConfigInstance.RedisAddr,
	Password: config.ConfigInstance.RedisPassword,
	DB: config.ConfigInstance.RedisDatabase,
})

func HGet(key, field string)string{
	rawStr, err := RedisClient.HGet(key, field).Result()
	if err != nil{
		utils.Loginfo("hget_error", utils.RedisFile, map[string]interface{}{
			"error": err,
			"key": key,
			"field": field,
		})
		return ""
	}
	return rawStr
}

func HSet(key string, field string, val interface{})error{
	r, err := RedisClient.HSet(key, field, val).Result()
	if err != nil{
		utils.Loginfo("hset_error", utils.RedisFile, map[string]interface{}{
			"result": r,
			"error": err,
			"key": key, "field": field, "value": val,
		})
		return err
	}
	return nil
}

func HDel(key string, field string){
	r, err := RedisClient.HDel(key, field).Result()
	if err != nil{
		utils.Loginfo("hdel_error", utils.RedisFile, map[string]interface{}{
			"error": err, "result": r,
			"key": key, "field": field,
		})
	}
}

func HGetAll(key string)map[string]string{
	result, err := RedisClient.HGetAll(key).Result()
	if err != nil{
		utils.Loginfo("hdel_error", utils.RedisFile, map[string]interface{}{
			"error": err, "result": result,
			"key": key,
		})
		return nil
	}
	return result
}