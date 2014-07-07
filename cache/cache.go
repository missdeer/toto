package cache

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/garyburd/redigo/redis"
	"github.com/missdeer/KellyBackend/setting"
)

var Mc *memcache.Client
var Rd redis.Conn

func init() {
	if setting.MemcachedEnabled {
		Mc = memcache.New(setting.MemcachedConn)
	}
	if setting.RedisEnabled {
		Rd, _ = redis.Dial("tcp", setting.RedisConn)
	}
	fmt.Println("initialize post package")
}
