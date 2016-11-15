package heartwater

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

const (
	TYPE_FOOTBALL = iota
	TYPE_BASKETBALL
)

type HeartwaterRouter struct {
	base.BaseRouter
}

func (this *HeartwaterRouter) Basketball() {
	if this.CheckLoginRedirect() {
		return
	}
	this.Data["IsHome"] = false
	// read from memcached or redis

	var res []models.HeartwaterRecord
	var err error
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetHeartwater("hw-basket", &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetHeartwater("hw-basket", &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Heartwater"] = res
		this.Data["RecordNum"] = len(res)
	}
	this.Data["Type"] = TYPE_BASKETBALL

	this.TplName = "heartwater/heartwater.html"
}

func (this *HeartwaterRouter) FootballLeague() {
	if this.CheckLoginRedirect() {
		return
	}
	this.Data["IsHome"] = false
	// read from memcached or redis

	var res []models.HeartwaterRecord
	var err error
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetHeartwater("hw-football", &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetHeartwater("hw-football", &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		// filter the league
		leagueId := this.GetString(":id")
		var r []models.HeartwaterRecord
		for _, v := range res {
			if v.LeagueId == leagueId {
				r = append(r, v)
			}
		}

		this.Data["Heartwater"] = r
		this.Data["RecordNum"] = len(r)
	}
	this.Data["Type"] = TYPE_FOOTBALL

	this.TplName = "heartwater/heartwater.html"
}

func (this *HeartwaterRouter) Football() {
	if this.CheckLoginRedirect() {
		return
	}
	this.Data["IsHome"] = false
	// read from memcached or redis

	var res []models.HeartwaterRecord
	var err error
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetHeartwater("hw-football", &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetHeartwater("hw-football", &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Heartwater"] = res
		this.Data["RecordNum"] = len(res)
	}
	this.Data["Type"] = TYPE_FOOTBALL

	this.TplName = "heartwater/heartwater.html"
}

func (this *HeartwaterRouter) fetchFootballDataSource() error {
	// read from data source and save to memcached or redis
	url := "http://zqcf2010.com:8080/front/recommend/game/list"
	resp, err := http.Get(url)
	if err != nil {
		beego.Error("read response from http://zqcf2010.com:8080/front/recommend/game/list error: ", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("read body from heartwater response error: ", err)
		return err
	}
	beginPos := bytes.Index(body, []byte("[{"))
	if beginPos == -1 {
		beego.Error("unexpect response: ", string(body))
		if string(body) == `var gameList=[];` {
			// clear memcached & redis
			if setting.MemcachedEnabled {
				cache.MemcachedRemove("hw-football")
			}

			if setting.RedisEnabled {
				cache.RedisRemove("hw-football")
			}
		}
		return err
	}
	body = body[beginPos : len(body)-1]

	var res []models.HeartwaterRecord
	if err = json.Unmarshal(body, &res); err != nil {
		beego.Error("json unmarshalling data source failed: ", err)
		return err
	}

	var result []models.HeartwaterRecord
	for _, r := range res {
		t, err := time.Parse("2006-01-02 15:04 -0700", r.StartTime+" +0800")
		if err != nil {
			beego.Error("parse time failed: ", err)
			continue
		}
		if t.After(time.Now()) {
			result = append(result, r)
		}
	}

	if setting.MemcachedEnabled {
		cache.MemcachedSetHeartwater("hw-football", &result)
	}

	if setting.RedisEnabled {
		cache.RedisSetHeartwater("hw-football", &result)
	}

	return nil
}

func (this *HeartwaterRouter) FetchFromDataSource() {
	this.fetchFootballDataSource()
	timer := time.NewTicker(15 * time.Second) // update data every 15 seconds
	for {
		select {
		case <-timer.C:
			this.fetchFootballDataSource()
		}
	}
	timer.Stop()
}
