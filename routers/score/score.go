package score

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/missdeer/toto/cache"
	"github.com/missdeer/toto/modules/models"
	"github.com/missdeer/toto/routers/base"
	"github.com/missdeer/toto/setting"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	teamMap = map[int]string{
		8:   `英超`,
		70:  `英冠`,
		13:  `意甲`,
		7:   `西甲`,
		9:   `德甲`,
		51:  `中超`,
		16:  `法甲`,
		17:  `法乙`,
		1:   `荷甲`,
		63:  `葡超`,
		121: `俄超`,
		26:  `巴甲`,
	}
)

type ScoreRouter struct {
	base.BaseRouter
}

func (this *ScoreRouter) FootballShooters() {
	// read from memcached or redis

	var res models.FootballScore
	id, err := this.GetInt(":id")
	if err != nil {
		id = 13
	}
	key := fmt.Sprintf("score-fb%d", id)
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetFootballScore(key, &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetFootballScore(key, &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Score"] = res.Playerrank[0].Playerrank
		this.Data["Team"] = teamMap
		this.Data["TeamId"] = id
	}

	this.TplNames = "score/shooters.html"
}

func (this *ScoreRouter) FootballAssistants() {
	// read from memcached or redis

	var res models.FootballScore
	id, err := this.GetInt(":id")
	if err != nil {
		id = 13
	}

	key := fmt.Sprintf("score-fb%d", id)
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetFootballScore(key, &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetFootballScore(key, &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Score"] = res.Assistrank[0].Playerassistrank
		this.Data["Team"] = teamMap
		this.Data["TeamId"] = id
	}

	this.TplNames = "score/assistants.html"
}

func (this *ScoreRouter) FootballCards() {
	// read from memcached or redis

	var res models.FootballScore
	id, err := this.GetInt(":id")
	if err != nil {
		id = 13
	}

	key := fmt.Sprintf("score-fb%d", id)
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetFootballScore(key, &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetFootballScore(key, &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Score"] = res.Cardrank[0].Standings
		this.Data["Team"] = teamMap
		this.Data["TeamId"] = id
	}

	this.TplNames = "score/cards.html"
}

func (this *ScoreRouter) FootballStandings() {
	// read from memcached or redis

	var res models.FootballScore
	id, err := this.GetInt(":id")
	if err != nil {
		id = 13
	}

	key := fmt.Sprintf("score-fb%d", id)
	if setting.MemcachedEnabled {
		err = cache.MemcachedGetFootballScore(key, &res)
	}

	if setting.RedisEnabled {
		err = cache.RedisGetFootballScore(key, &res)
	}

	if err != nil {
		this.Data["RecordNum"] = 0
	} else {
		this.Data["Score"] = res.Standings[0].Standings
		this.Data["Team"] = teamMap
		this.Data["TeamId"] = id
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

	var res models.FootballScore
	if err = json.Unmarshal(body, &res); err != nil {
		beego.Error("json unmarshalling data source failed: ", err)
		return err
	}

	key := fmt.Sprintf("score-fb%d", id)
	if setting.MemcachedEnabled {
		cache.MemcachedSetFootballScore(key, &res)
	}

	if setting.RedisEnabled {
		cache.RedisSetFootballScore(key, &res)
	}

	return nil
}

func (this *ScoreRouter) FetchFromDataSource() {
	for id, _ := range teamMap {
		this.fetchFootballDataSource(id)
	}

	timer := time.NewTicker(1 * time.Hour) // update data every 1hour
	for {
		select {
		case <-timer.C:
			for id, _ := range teamMap {
				this.fetchFootballDataSource(id)
			}
		}
	}
	timer.Stop()
}
