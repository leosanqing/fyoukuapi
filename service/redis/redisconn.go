package redisconn

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	ZAdd       = "ZADD"
	Exists     = "EXISTS"
	HmSet      = "HMSET"
	HGetAll    = "HGETALL"
	Expire     = "EXPIRE"
	WithScores = "WITHSCORES"
	ZRevRange  = "ZREVRANGE"
	LLen       = "LLEN"
	LRange     = "LRANGE"
	RPush      = "RPUSH"
)

func PoolConnect() redis.Conn {
	pool := &redis.Pool{
		// 最大空闲连接数
		MaxIdle: 1,
		// 最大连接数
		MaxActive: 10,
		// 超时时间
		IdleTimeout: 180 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", beego.AppConfig.String("redisdb"))
			if err != nil {
				return nil, err
			}
			return conn, err
		},
	}
	return pool.Get()
}
