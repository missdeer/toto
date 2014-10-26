package cache

import (
	"github.com/astaxie/beego"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/garyburd/redigo/redis"
	"github.com/missdeer/toto/setting"
)

var Mc *memcache.Client
var Rd redis.Conn

func Init() {
	if setting.MemcachedEnabled {
		beego.Info("enabled memcached at " + setting.MemcachedConn)
		Mc = memcache.New(setting.MemcachedConn)
	} else {
		beego.Warn("memcached is not enabled")
	}

	if setting.RedisEnabled {
		beego.Info("enabled redis at " + setting.RedisConn)
		Rd, _ = redis.Dial("tcp", setting.RedisConn)
	} else {
		beego.Warn("Redis is not enabled")
	}
}
