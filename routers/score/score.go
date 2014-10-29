package score

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

type ScoreRouter struct {
	base.BaseRouter
}

func (this *ScoreRouter) FootballShooters() {
	// read from memcached or redis

	var res []models.ScoreRecord
	var err error
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetScore("score-fb", &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetScore("score-fb", &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Score"] = res
		this.Data["RecordNum"] = len(res)
	}

	this.TplNames = "score/shooters.html"
}

func (this *ScoreRouter) FootballAssistants() {
	// read from memcached or redis

	var res []models.ScoreRecord
	var err error
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetScore("score-fb", &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetScore("score-fb", &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Score"] = res
		this.Data["RecordNum"] = len(res)
	}

	this.TplNames = "score/assistants.html"
}

func (this *ScoreRouter) FootballCards() {
	// read from memcached or redis

	var res []models.ScoreRecord
	var err error
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetScore("score-fb", &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetScore("score-fb", &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Score"] = res
		this.Data["RecordNum"] = len(res)
	}

	this.TplNames = "score/cards.html"
}

func (this *ScoreRouter) FootballStandings() {
	// read from memcached or redis

	var res []models.ScoreRecord
	var err error
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetScore("score-fb", &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetScore("score-fb", &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Score"] = res
		this.Data["RecordNum"] = len(res)
	}

	this.TplNames = "score/standings.html"
}

func (this *ScoreRouter) fetchFootballDataSource(id int) error {
	// read from data source and save to memcached or redis
	url := fmt.Sprintf("http://www.dongqiudi.com/getjson/%d?url=%2Fgetjson%2F%d", id, id)
	resp, err := http.Get(url)
	if err != nil {
		beego.Error("read response from ", url, err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("read body from score response error: ", err)
		return err
	}

	var res []models.ScoreRecord
	if err = json.Unmarshal(body, &res); err != nil {
		beego.Error("json unmarshalling data source failed: ", err)
		return err
	}

	if setting.MemcachedEnabled {
		cache.MemcachedSetScore("score-fb", &res)
	}

	if setting.RedisEnabled {
		cache.RedisSetScore("score-fb", &res)
	}

	return nil
}

func (this *ScoreRouter) FetchFromDataSource() {
	ids := []int{8, 70, 13, 7, 9, 51, 16, 17, 1, 63, 121, 26}
	for _, id := range ids {
		this.fetchFootballDataSource(id)
	}

	timer := time.NewTicker(1 * time.Hour) // update data every 1hour
	for {
		select {
		case <-timer.C:
			for _, id := range ids {
				this.fetchFootballDataSource(id)
			}
		}
	}
	timer.Stop()
}
