package config

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload" // 引入.env 变量
)

// redis的缓存连接信息
var RedisParseURL string

// mysql的数据库链接信息
var MysqlParseURL string

// 暂未对外提供，所以留存密钥
var JwtSecret string = "160ee67dd3ecc449fc526f85"

func init() {
	gin.SetMode(os.Getenv("GIN_MODE"))
	RedisParseURL = "redis://" + os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT") + "/" + os.Getenv("REDIS_DB")
	MysqlParseURL = os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
}
