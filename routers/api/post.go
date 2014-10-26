// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package api

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"

	"github.com/missdeer/toto/cache"
	"github.com/missdeer/toto/modules/models"
	"github.com/missdeer/toto/setting"
)

func (this *ApiRouter) PostToggle() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()

	if !this.IsAjax() {
		return
	}

	action := this.GetString("action")

	if this.IsLogin {

		switch action {
		case "toggle-best":
			id, _ := this.GetInt("post")
			if id > 0 {
				o := orm.NewOrm()
				p := models.Post{Id: int(id)}
				o.Read(&p)

				p.IsBest = !p.IsBest
				if _, err := o.Update(&p); err != nil {
					beego.Error("PostCounterAdd ", err)
				} else {
					result["success"] = true
					// update home/recent/category/topic/best posts cache
					if setting.MemcachedEnabled {
						cache.Mc.Delete("recent-posts-count")
						cache.Mc.Delete("recent-posts")
						cache.Mc.Delete("recent-category")
						cache.Mc.Delete("home-posts")
						cache.Mc.Delete("best-posts-count")
						cache.Mc.Delete("best-posts")
						cache.Mc.Delete("best-category")
						categoryCountKey := fmt.Sprintf(`category-%s-count`, p.Category.Slug)
						cache.Mc.Delete(categoryCountKey)
						categoryKey := fmt.Sprintf(`category-%s`, p.Category.Slug)
						cache.Mc.Delete(categoryKey)
						topicCountKey := fmt.Sprintf(`topic-%s-count`, p.Topic.Slug)
						cache.Mc.Delete(topicCountKey)
						topicKey := fmt.Sprintf(`topic-%s`, p.Topic.Slug)
						cache.Mc.Delete(topicKey)
					}
					if setting.RedisEnabled {
						cache.Rd.Do("DEL", "recent-posts-count")
						cache.Rd.Do("DEL", "recent-posts")
						cache.Rd.Do("DEL", "recent-category")
						cache.Rd.Do("DEL", "home-posts")
						cache.Rd.Do("DEL", "best-posts-count")
						cache.Rd.Do("DEL", "best-posts")
						cache.Rd.Do("DEL", "best-category")
						categoryCountKey := fmt.Sprintf(`category-%s-count`, p.Category.Slug)
						cache.Rd.Do("DEL", categoryCountKey)
						categoryKey := fmt.Sprintf(`category-%s`, p.Category.Slug)
						cache.Rd.Do("DEL", categoryKey)
						topicCountKey := fmt.Sprintf(`topic-%s-count`, p.Topic.Slug)
						cache.Rd.Do("DEL", topicCountKey)
						topicKey := fmt.Sprintf(`topic-%s`, p.Topic.Slug)
						cache.Rd.Do("DEL", topicKey)
					}
				}
				o = nil
			}
		case "toggle-top":
			id, _ := this.GetInt("post")
			if id > 0 {
				o := orm.NewOrm()
				p := models.Post{Id: int(id)}
				o.Read(&p)

				p.IsTop = !p.IsTop
				if _, err := o.Update(&p); err != nil {
					beego.Error("PostCounterAdd ", err)
				} else {
					result["success"] = true
					// update home/recent/category/topic posts cache
					if setting.MemcachedEnabled {
						cache.Mc.Delete("recent-posts-count")
						cache.Mc.Delete("recent-posts")
						cache.Mc.Delete("recent-category")
						cache.Mc.Delete("home-posts")
						categoryCountKey := fmt.Sprintf(`category-%s-count`, p.Category.Slug)
						cache.Mc.Delete(categoryCountKey)
						categoryKey := fmt.Sprintf(`category-%s`, p.Category.Slug)
						cache.Mc.Delete(categoryKey)
						topicCountKey := fmt.Sprintf(`topic-%s-count`, p.Topic.Slug)
						cache.Mc.Delete(topicCountKey)
						topicKey := fmt.Sprintf(`topic-%s`, p.Topic.Slug)
						cache.Mc.Delete(topicKey)
					}
					if setting.RedisEnabled {
						cache.Rd.Do("DEL", "recent-posts-count")
						cache.Rd.Do("DEL", "recent-posts")
						cache.Rd.Do("DEL", "recent-category")
						cache.Rd.Do("DEL", "home-posts")
						categoryCountKey := fmt.Sprintf(`category-%s-count`, p.Category.Slug)
						cache.Rd.Do("DEL", categoryCountKey)
						categoryKey := fmt.Sprintf(`category-%s`, p.Category.Slug)
						cache.Rd.Do("DEL", categoryKey)
						topicCountKey := fmt.Sprintf(`topic-%s-count`, p.Topic.Slug)
						cache.Rd.Do("DEL", topicCountKey)
						topicKey := fmt.Sprintf(`topic-%s`, p.Topic.Slug)
						cache.Rd.Do("DEL", topicKey)
					}
				}
				o = nil
			}
		}
	}
}

func ClearTodayReplys() {
	timer := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-timer.C:
			now := time.Now().UTC()
			if now.Hour() == 16 && now.Minute() == 0 {
				// clear it when it's 00:00 at GMT+8 (Asia/Shanghai) time zone
				beego.Info("clear today replys of all posts")
				o := orm.NewOrm()
				_, err := o.QueryTable("post").Update(orm.Params{
					"today_replys": 0,
				})
				if err != nil {
					beego.Error("clear today replys error ", err)
				}
				o = nil
			}
		}
	}
	timer.Stop()
}
