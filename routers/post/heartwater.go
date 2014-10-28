package post

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/missdeer/toto/cache"
	"github.com/missdeer/toto/modules/models"
	"github.com/missdeer/toto/routers/base"
	"github.com/missdeer/toto/setting"
	"io/ioutil"
	"net/http"
	"time"
)

type HeartwaterRouter struct {
	base.BaseRouter
}

func (this *HeartwaterRouter) Get() {
	if this.CheckLoginRedirect() {
		return
	}
	this.Data["IsHome"] = false
	this.TplNames = "post/heartwater.html"
}

func (this *HeartwaterRouter) FetchFromDataSource() {
	timer := time.NewTicker(15 * time.Second) // update data every 15 seconds
	for {
		select {
		case <-timer.C:
			// read from data source and save to memcached or redis
			url := "http://zqcf2010.com:8080/front/recommend/game/list"
			resp, err := http.Get(url)
			if err != nil {
				beego.Error("read response from http://zqcf2010.com:8080/front/recommend/game/list error: ", err)
				break
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				beego.Error("read body from heartwater response error: ", err)
				break
			}
			beginPos := bytes.Index(body, []byte("[{"))
			body = body[beginPos : len(body)-1]

			var res []models.HeartwaterRecord
			if err = json.Unmarshal(body, &res); err != nil {
				beego.Error("json unmarshalling data source failed: ", err)
				break
			}

			if setting.MemcachedEnabled {
				cache.MemcachedSetHeartwater("heartwater", &res)
			}

			if setting.RedisEnabled {
				cache.RedisSetHeartwater("heartwater", &res)
			}
		}
	}
	timer.Stop()
}
