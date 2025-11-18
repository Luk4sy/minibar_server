package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/redis_service/redis_jwt"
	"blogx_server/utils/jwts"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Redis = core.InitRedis()

	token, err := jwts.GetToken(jwts.Claims{
		UserID: 2,
		Role:   1,
	})
	fmt.Println(token, err)
	redis_jwt.TokenBlackList(token, redis_jwt.UserBlackListType)
	blk, ok := redis_jwt.HasTokenBlackList(token)
	fmt.Println(blk, ok)
}
