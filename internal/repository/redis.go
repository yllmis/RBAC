package repository

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var store *redisstore.RedisStore
var sessionName = "session-name"

func NewRedis() {
	//redis连接
	fmt.Println("正在连接Redis数据库......")

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	Rdb = rdb
	store, _ = redisstore.NewRedisStore(context.TODO(), Rdb)
	//redis连接成功
	fmt.Println("Redis数据库连接成功!")
}

func GetSession(ctx *gin.Context) map[interface{}]interface{} {
	session, err := store.Get(ctx.Request, sessionName)
	fmt.Printf("session: %+v, err: %+v\n", session, err)
	return session.Values
}

func SetSession(ctx gin.Context, name string, id int64) error {
	session, err := store.Get(ctx.Request, sessionName)
	if err != nil {
		fmt.Printf("获取session失败: %+v\n", err)
		return err
	}
	session.Values[name] = name
	session.Values["id"] = id
	return session.Save(ctx.Request, ctx.Writer)
}

func FlushSession(ctx gin.Context) error {
	session, err := store.Get(ctx.Request, sessionName)
	if err != nil {
		fmt.Printf("获取session失败: %+v\n", err)
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(ctx.Request, ctx.Writer)
}
