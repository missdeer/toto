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

package post

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/missdeer/toto/cache"
	"github.com/missdeer/toto/modules/models"
	"github.com/missdeer/toto/modules/post"
	"github.com/missdeer/toto/modules/utils"
	"github.com/missdeer/toto/routers/base"
	"github.com/missdeer/toto/setting"
)

// HomeRouter serves home page.
type PostListRouter struct {
	base.BaseRouter
}

func (this *PostListRouter) setCategories(cats *[]models.Category) {
	post.ListCategories(cats)
	this.Data["Categories"] = *cats
}

func (this *PostListRouter) setTopicsOfCat(topics *[]models.Topic, cat *models.Category) {
	post.ListTopicsOfCat(topics, cat)
	this.Data["Topics"] = *topics
}

func (this *PostListRouter) postsFilter(qs orm.QuerySeter) orm.QuerySeter {
	if !this.IsLogin {
		return qs
	}
	args := []string{utils.ToStr(this.Locale.Index())}
	args = append(args, this.User.LangAdds...)
	args = append(args, utils.ToStr(this.User.Lang))
	qs = qs.Filter("Lang__in", args)
	return qs
}

func (this *PostListRouter) ORCA() {
	orca_verify_code := setting.ORCAVerifyCode
	this.Ctx.WriteString(orca_verify_code)
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) Home() {
	this.Data["IsHome"] = true
	this.TplNames = "post/home.html"
}

func (this *PostListRouter) BBS() {
	this.Data["IsHome"] = false
	this.TplNames = "post/bbs.html"

	var cats []models.Category
	this.setCategories(&cats)

	this.Data["CategorySlug"] = "hot"

	// get topics
	var topics []models.Topic
	if setting.MemcachedEnabled {
		if err := cache.MemcachedGetTopics("home-topics", &topics); err != nil {
			beego.Error("get home topics from memcache failed. ", err)
			post.ListTopics(&topics)
			cache.MemcachedSetTopics("home-topics", &topics)
		}
	}

	if setting.RedisEnabled {
		if err := cache.RedisGetTopics("home-topics", &topics); err != nil {
			beego.Error("get home topics from redis failed. ", err)
			post.ListTopics(&topics)
			cache.RedisSetTopics("home-topics", &topics)
		}
	}

	if len(topics) == 0 {
		post.ListTopics(&topics)
	}
	this.Data["Topics"] = topics

	// get posts
	var posts []models.Post
	var todayTopTen []models.Post
	if setting.MemcachedEnabled {
		if cache.MemcachedGetPosts("home-posts", &posts) == nil && cache.MemcachedGetPosts("today-topten-posts", &todayTopTen) == nil {
			this.Data["Posts"] = posts
			this.Data["TodayTopTen"] = todayTopTen
			return
		}
	}

	if setting.RedisEnabled {
		if cache.RedisGetPosts("home-posts", &posts) == nil && cache.RedisGetPosts("today-topten-posts", &todayTopTen) == nil {
			this.Data["Posts"] = posts
			this.Data["TodayTopTen"] = todayTopTen
			return
		}
	}

	beego.Info("Get posts from database")
	postsModel := models.Posts()

	var topposts []models.Post
	qs := postsModel.Exclude("category_id", setting.CategoryHideOnHome).Filter("IsTop", true).OrderBy("-Created").Limit(25).RelatedSel()
	qs = this.postsFilter(qs)
	models.ListObjects(qs, &topposts)

	topCount := len(topposts)
	qs2 := postsModel.Exclude("category_id", setting.CategoryHideOnHome).Filter("IsTop", false).OrderBy("-Created").Limit(25 - topCount).RelatedSel()
	qs2 = this.postsFilter(qs2)
	var nontopposts []models.Post
	models.ListObjects(qs2, &nontopposts)

	posts = append(topposts, nontopposts...)

	qsTopTen := postsModel.Exclude("today_replys", 0).OrderBy("-TodayReplys").Limit(10).RelatedSel()
	qsTopTen = this.postsFilter(qsTopTen)
	models.ListObjects(qsTopTen, &todayTopTen)

	this.Data["Posts"] = posts
	this.Data["TodayTopTen"] = todayTopTen

	if setting.MemcachedEnabled {
		cache.MemcachedSetPosts("home-posts", &posts)
		cache.MemcachedSetPosts("today-topten-posts", &todayTopTen)
	}

	if setting.RedisEnabled {
		cache.RedisSetPosts("home-posts", &posts)
		cache.RedisSetPosts("today-topten-posts", &todayTopTen)
	}
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) Category() {
	this.TplNames = "post/category.html"

	slug := this.GetString(":slug")
	cat := models.Category{Slug: slug}
	if err := cat.Read("Slug"); err != nil {
		this.Abort("404")
		return
	}

	this.Data["Category"] = &cat
	this.Data["CategorySlug"] = cat.Slug
	this.Data["IsCategory"] = true

	var cats []models.Category
	this.setCategories(&cats)

	var topics []models.Topic
	this.setTopicsOfCat(&topics, &cat)

	var posts []models.Post
	pers := 25
	var cnt int64
	var pager *utils.Paginator
	var err error

	if setting.MemcachedEnabled {
		key := fmt.Sprintf("category-%s-count", slug)
		if cnt, err = cache.MemcachedGetInt64(key); err == nil {
			pager = this.SetPaginator(pers, cnt)
			if pager.Page() == 1 {
				key = fmt.Sprintf("category-%s", slug)
				if cache.MemcachedGetPosts(key, &posts) == nil {
					this.Data["Posts"] = posts
					return
				}
			}
		}
	}

	if setting.RedisEnabled {
		key := fmt.Sprintf("category-%s-count", slug)
		if cnt, err = cache.RedisGetInt64(key); err == nil {
			pager = this.SetPaginator(pers, cnt)
			if pager.Page() == 1 {
				key = fmt.Sprintf("category-%s", slug)
				if cache.RedisGetPosts(key, &posts) == nil {
					this.Data["Posts"] = posts
					return
				}
			}
		}
	}

	qs := models.Posts().Filter("Category", &cat)
	qs = this.postsFilter(qs)

	cnt, _ = models.CountObjects(qs)
	pager = this.SetPaginator(pers, cnt)
	if setting.MemcachedEnabled {
		key := fmt.Sprintf("category-%s-count", slug)
		cache.MemcachedSetInt64(key, cnt)
	}

	if setting.RedisEnabled {
		key := fmt.Sprintf("category-%s-count", slug)
		cache.RedisSetInt64(key, cnt)
	}

	if pager.Page() > 1 {
		qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()
		models.ListObjects(qs, &posts)
	} else {
		qsTop := models.Posts().Filter("Category", &cat).Filter("IsTop", true)
		qsTop = this.postsFilter(qsTop).OrderBy("-Created").Limit(pers).RelatedSel()
		var topposts []models.Post
		models.ListObjects(qsTop, &topposts)

		qsNonTop := models.Posts().Filter("Category", &cat).Filter("IsTop", false)
		qsNonTop = this.postsFilter(qsNonTop).OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()
		var nontopposts []models.Post
		models.ListObjects(qsNonTop, &nontopposts)

		posts = append(topposts, nontopposts...)

		if setting.MemcachedEnabled {
			key := fmt.Sprintf("category-%s", slug)
			cache.MemcachedSetPosts(key, &posts)
		}

		if setting.RedisEnabled {
			key := fmt.Sprintf("category-%s", slug)
			cache.RedisSetPosts(key, &posts)
		}
	}

	this.Data["Posts"] = posts
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) Navs() {
	slug := this.GetString(":slug")

	switch slug {
	case "favs", "follow":
		if this.CheckLoginRedirect() {
			return
		}
	}

	this.Data["CategorySlug"] = slug
	this.TplNames = fmt.Sprintf("post/navs/%s.html", slug)

	pers := 25

	var posts []models.Post
	var cats []models.Category
	var cnt int64
	var pager *utils.Paginator
	var err error

	switch slug {
	case "recent":
		if setting.MemcachedEnabled {
			if cnt, err = cache.MemcachedGetInt64("recent-posts-count"); err == nil {
				pager = this.SetPaginator(pers, cnt)
				if pager.Page() == 1 {
					if cache.MemcachedGetPosts("recent-posts", &posts) == nil && cache.MemcachedGetCategories("recent-category", &cats) == nil {
						this.Data["Categories"] = cats
						break
					}
				}
			}
		}

		if setting.RedisEnabled {
			if cnt, err = cache.RedisGetInt64("recent-posts-count"); err == nil {
				pager = this.SetPaginator(pers, cnt)
				if pager.Page() == 1 {
					if cache.RedisGetPosts("recent-posts", &posts) == nil && cache.RedisGetCategories("recent-category", &cats) == nil {
						this.Data["Categories"] = cats
						break
					}
				}
			}
		}
		qs := models.Posts().Exclude("category_id", setting.CategoryHideOnHome)
		qs = this.postsFilter(qs)

		cnt, _ = models.CountObjects(qs)
		pager = this.SetPaginator(pers, cnt)

		qs = qs.OrderBy("-Updated").Limit(pers, pager.Offset()).RelatedSel()

		models.ListObjects(qs, &posts)

		this.setCategories(&cats)

		if setting.MemcachedEnabled {
			cache.MemcachedSetInt64("recent-posts-count", cnt)
			cache.MemcachedSetPosts("recent-posts", &posts)
			cache.MemcachedSetCategories("recent-category", &cats)
		}

		if setting.RedisEnabled {
			cache.RedisSetInt64("recent-posts-count", cnt)
			cache.RedisSetPosts("recent-posts", &posts)
			cache.RedisSetCategories("recent-category", &cats)
		}
	case "best":
		if setting.MemcachedEnabled {
			if cnt, err = cache.MemcachedGetInt64("best-posts-count"); err == nil {
				pager = this.SetPaginator(pers, cnt)
				if pager.Page() == 1 {
					if cache.MemcachedGetPosts("best-posts", &posts) == nil && cache.MemcachedGetCategories("best-category", &cats) == nil {
						this.Data["Categories"] = cats
						break
					}
				}
			}
		}

		if setting.RedisEnabled {
			if cnt, err = cache.RedisGetInt64("best-posts-count"); err == nil {
				pager = this.SetPaginator(pers, cnt)
				if pager.Page() == 1 {
					if cache.RedisGetPosts("best-posts", &posts) == nil && cache.RedisGetCategories("best-category", &cats) == nil {
						this.Data["Categories"] = cats
						break
					}
				}
			}
		}

		qs := models.Posts().Filter("IsBest", true)
		qs = this.postsFilter(qs)

		cnt, _ := models.CountObjects(qs)
		pager := this.SetPaginator(pers, cnt)

		qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

		models.ListObjects(qs, &posts)
		this.setCategories(&cats)

		if setting.MemcachedEnabled {
			cache.MemcachedSetInt64("best-posts-count", cnt)
			cache.MemcachedSetPosts("best-posts", &posts)
			cache.MemcachedSetCategories("best-category", &cats)
		}

		if setting.RedisEnabled {
			cache.RedisSetInt64("best-posts-count", cnt)
			cache.RedisSetPosts("best-posts", &posts)
			cache.RedisSetCategories("best-category", &cats)
		}
	case "cold":
		if setting.MemcachedEnabled {
			if cnt, err = cache.MemcachedGetInt64("cold-posts-count"); err == nil {
				pager = this.SetPaginator(pers, cnt)
				if pager.Page() == 1 {
					if cache.MemcachedGetPosts("cold-posts", &posts) == nil && cache.MemcachedGetCategories("cold-category", &cats) == nil {
						this.Data["Categories"] = cats
						break
					}
				}
			}
		}

		if setting.RedisEnabled {
			if cnt, err = cache.RedisGetInt64("cold-posts-count"); err == nil {
				pager = this.SetPaginator(pers, cnt)
				if pager.Page() == 1 {
					if cache.RedisGetPosts("cold-posts", &posts) == nil && cache.RedisGetCategories("cold-category", &cats) == nil {
						this.Data["Categories"] = cats
						break
					}
				}
			}
		}
		qs := models.Posts().Filter("Replys", 0)
		qs = this.postsFilter(qs)

		cnt, _ := models.CountObjects(qs)
		pager := this.SetPaginator(pers, cnt)

		qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

		models.ListObjects(qs, &posts)
		this.setCategories(&cats)

		if setting.MemcachedEnabled {
			cache.MemcachedSetInt64("cold-posts-count", cnt)
			cache.MemcachedSetPosts("cold-posts", &posts)
			cache.MemcachedSetCategories("cold-category", &cats)
		}

		if setting.RedisEnabled {
			cache.RedisSetInt64("cold-posts-count", cnt)
			cache.RedisSetPosts("cold-posts", &posts)
			cache.RedisSetCategories("cold-category", &cats)
		}
	case "favs":
		var topicIds orm.ParamsList
		nums, _ := models.FollowTopics().Filter("User", &this.User.Id).OrderBy("-Created").ValuesFlat(&topicIds, "Topic")
		if nums > 0 {
			qs := models.Posts().Filter("Topic__in", topicIds)
			qs = this.postsFilter(qs)

			cnt, _ := models.CountObjects(qs)
			pager := this.SetPaginator(pers, cnt)

			qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

			models.ListObjects(qs, &posts)

			var topics []models.Topic
			nums, _ = models.Topics().Filter("Id__in", topicIds).Limit(8).All(&topics)
			this.Data["Topics"] = topics
			this.Data["TopicsMore"] = nums >= 8
		}

	case "follow":
		var userIds orm.ParamsList
		nums, _ := this.User.FollowingUsers().OrderBy("-Created").ValuesFlat(&userIds, "FollowUser")
		if nums > 0 {
			qs := models.Posts().Filter("User__in", userIds)
			qs = this.postsFilter(qs)

			cnt, _ := models.CountObjects(qs)
			pager := this.SetPaginator(pers, cnt)

			qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

			models.ListObjects(qs, &posts)
		}
	}

	this.Data["Posts"] = posts
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) Topic() {
	slug := this.GetString(":slug")

	switch slug {
	default: // View topic.
		this.TplNames = "post/topic.html"
		topic := models.Topic{Slug: slug}
		if err := topic.Read("Slug"); err != nil {
			this.Abort("404")
			return
		}

		this.Data["Slug"] = slug
		this.Data["Topic"] = &topic
		this.Data["IsTopic"] = true

		HasFavorite := false
		if this.IsLogin {
			HasFavorite = models.FollowTopics().Filter("User", &this.User).Filter("Topic", &topic).Exist()
		}
		this.Data["HasFavorite"] = HasFavorite

		var posts []models.Post
		pers := 25
		var cnt int64
		var pager *utils.Paginator
		var err error

		if setting.MemcachedEnabled {
			key := fmt.Sprintf("topic-%s-count", slug)
			if cnt, err = cache.MemcachedGetInt64(key); err == nil {
				pager = this.SetPaginator(pers, cnt)
				if pager.Page() == 1 {
					key = fmt.Sprintf("topic-%s", slug)
					if cache.MemcachedGetPosts(key, &posts) == nil {
						this.Data["Posts"] = posts
						return
					}
				}
			}
		}

		if setting.RedisEnabled {
			key := fmt.Sprintf("topic-%s-count", slug)
			if cnt, err = cache.RedisGetInt64(key); err == nil {
				pager = this.SetPaginator(pers, cnt)
				if pager.Page() == 1 {
					key = fmt.Sprintf("topic-%s", slug)
					if cache.RedisGetPosts(key, &posts) == nil {
						this.Data["Posts"] = posts
						return
					}
				}
			}
		}
		qs := models.Posts().Filter("Topic", &topic)
		qs = this.postsFilter(qs)

		cnt, _ = models.CountObjects(qs)
		pager = this.SetPaginator(pers, cnt)
		if setting.MemcachedEnabled {
			key := fmt.Sprintf("topic-%s-count", slug)
			cache.MemcachedSetInt64(key, cnt)
		}

		if setting.RedisEnabled {
			key := fmt.Sprintf("topic-%s-count", slug)
			cache.RedisSetInt64(key, cnt)
		}

		if pager.Page() > 1 {
			qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()
			models.ListObjects(qs, &posts)
		} else {
			qsTop := models.Posts().Filter("Topic", &topic).Filter("IsTop", true)
			qsTop = this.postsFilter(qsTop).OrderBy("-Created").Limit(pers).RelatedSel()
			var topposts []models.Post
			models.ListObjects(qsTop, &topposts)

			qsNonTop := models.Posts().Filter("Topic", &topic).Filter("IsTop", false)
			qsNonTop = this.postsFilter(qsNonTop).OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()
			var nontopposts []models.Post
			models.ListObjects(qsNonTop, &nontopposts)

			posts = append(topposts, nontopposts...)

			if setting.MemcachedEnabled {
				key := fmt.Sprintf("topic-%s", slug)
				cache.MemcachedSetPosts(key, &posts)
			}
			if setting.RedisEnabled {
				key := fmt.Sprintf("topic-%s", slug)
				cache.RedisSetPosts(key, &posts)
			}
		}

		this.Data["Posts"] = posts
	}
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) TopicSubmit() {
	slug := this.GetString(":slug")

	topic := models.Topic{Slug: slug}
	if err := topic.Read("Slug"); err != nil {
		this.Abort("404")
		return
	}

	result := map[string]interface{}{
		"success": false,
	}

	if this.IsAjax() {
		action := this.GetString("action")
		switch action {
		case "favorite":
			if this.IsLogin {
				qs := models.FollowTopics().Filter("User", &this.User).Filter("Topic", &topic)
				if qs.Exist() {
					qs.Delete()
				} else {
					fav := models.FollowTopic{User: &this.User, Topic: &topic}
					fav.Insert()
				}
				topic.RefreshFollowers()
				this.User.RefreshFavTopics()
				result["success"] = true
			}
		}
	}

	this.Data["json"] = result
	this.ServeJson()
}

type PostRouter struct {
	base.BaseRouter
}

func (this *PostRouter) New() {
	this.TplNames = "post/new.html"

	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}

	if v := this.Ctx.GetCookie("post_topic"); len(v) > 0 {
		form.Topic, _ = utils.StrTo(v).Int()
	}

	if v := this.Ctx.GetCookie("post_cat"); len(v) > 0 {
		form.Category, _ = utils.StrTo(v).Int()
	}

	if v := this.Ctx.GetCookie("post_lang"); len(v) > 0 {
		form.Lang, _ = utils.StrTo(v).Int()
	} else {
		form.Lang = this.Locale.Index()
	}

	slug := this.GetString("topic")
	if len(slug) > 0 {
		topic := models.Topic{Slug: slug}
		topic.Read("Slug")
		form.Topic = topic.Id
		this.Data["Topic"] = &topic
	}

	post.ListCategories(&form.Categories)
	post.ListTopics(&form.Topics)
	this.SetFormSets(&form)
}

func (this *PostRouter) NewSubmit() {
	this.TplNames = "post/new.html"

	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}
	slug := this.GetString("topic")
	if len(slug) > 0 {
		topic := models.Topic{Slug: slug}
		topic.Read("Slug")
		form.Topic = topic.Id
		this.Data["Topic"] = &topic
	}

	post.ListCategories(&form.Categories)
	post.ListTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	var post models.Post
	if err := form.SavePost(&post, &this.User); err == nil {

		this.Ctx.SetCookie("post_topic", utils.ToStr(form.Topic), 1<<31-1, "/")
		this.Ctx.SetCookie("post_cat", utils.ToStr(form.Category), 1<<31-1, "/")
		this.Ctx.SetCookie("post_lang", utils.ToStr(form.Lang), 1<<31-1, "/")

		this.JsStorage("deleteKey", "post/new")
		this.Redirect(post.Link(), 302)

		if len(post.Category.Slug) == 0 {
			cate := models.Category{Id: form.Category}
			cate.Read("Id")
			post.Category.Slug = cate.Slug
			beego.Info("category slug:", post.Category.Slug)
		}

		if len(post.Topic.Slug) == 0 {
			topic := models.Topic{Id: form.Topic}
			topic.Read("Id")
			post.Topic.Slug = topic.Slug
			beego.Info("topic slug:", post.Topic.Slug)
		}

		// update recent/home/category/topics posts cache
		if setting.MemcachedEnabled {
			cache.MemcachedRemove("recent-posts-count")
			cache.MemcachedRemove("recent-posts")
			cache.MemcachedRemove("recent-category")
			cache.MemcachedRemove("home-posts")
			cache.MemcachedRemove("today-topten-posts")

			if len(post.Category.Slug) == 0 {
				key := fmt.Sprintf(`category-slug-%d`, form.Category)
				if s, e := cache.MemcachedGetString(key); e != nil {
					cate := models.Category{Id: form.Category}
					cate.Read("Id")
					cache.MemcachedSetString(key, &cate.Slug)
					post.Category.Slug = cate.Slug
				} else {
					post.Category.Slug = s
				}
				beego.Info("category slug:", post.Category.Slug)
			}

			if len(post.Topic.Slug) == 0 {
				key := fmt.Sprintf(`topic-slug-%d`, form.Topic)
				if s, e := cache.MemcachedGetString(key); e != nil {
					topic := models.Topic{Id: form.Topic}
					topic.Read("Id")
					cache.MemcachedSetString(key, &topic.Slug)
					post.Topic.Slug = topic.Slug
				} else {
					post.Topic.Slug = s
				}
				beego.Info("topic slug:", post.Topic.Slug)
			}

			categoryCountKey := fmt.Sprintf(`category-%s-count`, post.Category.Slug)
			cache.MemcachedRemove(categoryCountKey)
			categoryKey := fmt.Sprintf(`category-%s`, post.Category.Slug)
			cache.MemcachedRemove(categoryKey)
			topicCountKey := fmt.Sprintf(`topic-%s-count`, post.Topic.Slug)
			cache.MemcachedRemove(topicCountKey)
			topicKey := fmt.Sprintf(`topic-%s`, post.Topic.Slug)
			cache.MemcachedRemove(topicKey)
		}

		if setting.RedisEnabled {
			cache.RedisRemove("recent-posts-count")
			cache.RedisRemove("recent-posts")
			cache.RedisRemove("recent-category")
			cache.RedisRemove("home-posts")
			cache.RedisRemove("today-topten-posts")

			if len(post.Category.Slug) == 0 {
				key := fmt.Sprintf(`category-slug-%d`, form.Category)
				if s, e := cache.RedisGetString(key); e != nil {
					cate := models.Category{Id: form.Category}
					cate.Read("Id")
					cache.RedisSetString(key, &cate.Slug)
					post.Category.Slug = cate.Slug
				} else {
					post.Category.Slug = s
				}
				beego.Info("category slug:", post.Category.Slug)
			}

			if len(post.Topic.Slug) == 0 {
				key := fmt.Sprintf(`topic-slug-%d`, form.Topic)
				if s, e := cache.RedisGetString(key); e != nil {
					topic := models.Topic{Id: form.Topic}
					topic.Read("Id")
					cache.RedisSetString(key, &topic.Slug)
					post.Topic.Slug = topic.Slug
				} else {
					post.Topic.Slug = s
				}
				beego.Info("topic slug:", post.Topic.Slug)
			}

			categoryCountKey := fmt.Sprintf(`category-%s-count`, post.Category.Slug)
			cache.RedisRemove(categoryCountKey)
			categoryKey := fmt.Sprintf(`category-%s`, post.Category.Slug)
			cache.RedisRemove(categoryKey)
			topicCountKey := fmt.Sprintf(`topic-%s-count`, post.Topic.Slug)
			cache.RedisRemove(topicCountKey)
			topicKey := fmt.Sprintf(`topic-%s`, post.Topic.Slug)
			cache.RedisRemove(topicKey)
		}
	}
}

func (this *PostRouter) loadPost(post *models.Post, user *models.User) bool {
	id, _ := this.GetInt(":post")
	if id > 0 {
		qs := models.Posts().Filter("Id", id)
		if user != nil {
			qs = qs.Filter("User", user.Id)
		}
		qs.RelatedSel(1).One(post)
	}

	if post.Id == 0 {
		this.Abort("404")
		return true
	}

	this.Data["Post"] = post

	return false
}

func (this *PostRouter) loadAppends(post *models.Post, appends *[]*models.AppendPost) {
	qs := post.Appends()
	if num, err := qs.OrderBy("Id").All(appends); err == nil {
		this.Data["Appends"] = *appends
		this.Data["AppendsNum"] = num
	}
}

func (this *PostRouter) loadComments(post *models.Post, comments *[]*models.Comment) {
	qs := post.Comments().Filter("Duplicated", false)
	if _, err := qs.RelatedSel("User").OrderBy("Id").All(comments); err == nil {
		this.Data["Comments"] = *comments
		this.Data["CommentsNum"] = post.Replys
	}
}

func (this *PostRouter) isDuplicatedComment(post *models.Post, message string) bool {
	qs := post.Comments().Filter("Message", message).RelatedSel()
	num, _ := qs.Count()
	return num > 0
}

func (this *PostRouter) Single() {
	this.TplNames = "post/post.html"

	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return
	}

	var comments []*models.Comment
	this.loadComments(&postMd, &comments)

	var appends []*models.AppendPost
	this.loadAppends(&postMd, &appends)

	form := post.CommentForm{}
	this.SetFormSets(&form)

	post.PostBrowsersAdd(this.User.Id, this.Ctx.Input.IP(), &postMd)
}

func (this *PostRouter) SingleSubmit() {
	this.TplNames = "post/post.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return
	}

	var redir bool

	defer func() {
		if !redir {
			var comments []*models.Comment
			this.loadComments(&postMd, &comments)
		}
	}()

	form := post.CommentForm{}
	if !this.ValidFormSets(&form) {
		return
	}

	comment := models.Comment{}
	comment.Duplicated = this.isDuplicatedComment(&postMd, form.Message)
	if err := form.SaveComment(&comment, &this.User, &postMd); err == nil {
		this.JsStorage("deleteKey", "post/comment")
		this.Redirect(postMd.Link(), 302)
		redir = true

		post.PostReplysCount(&postMd)
	}
	// update cold posts cache
	if postMd.Replys == 1 {
		if setting.MemcachedEnabled {
			cache.MemcachedRemove("cold-posts-count")
			cache.MemcachedRemove("cold-posts")
			cache.MemcachedRemove("cold-category")
		}
		if setting.RedisEnabled {
			cache.RedisRemove("cold-posts-count")
			cache.RedisRemove("cold-posts")
			cache.RedisRemove("cold-category")
		}
	}
}

func (this *PostRouter) Edit() {
	this.TplNames = "post/edit.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	form := post.PostForm{}
	form.SetFromPost(&postMd)
	post.ListCategories(&form.Categories)
	post.ListTopics(&form.Topics)
	this.SetFormSets(&form)
}

func (this *PostRouter) EditSubmit() {
	this.TplNames = "post/edit.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	form := post.PostForm{}
	form.SetFromPost(&postMd)
	post.ListCategories(&form.Categories)
	post.ListTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	if err := form.UpdatePost(&postMd, &this.User); err == nil {
		this.JsStorage("deleteKey", "post/edit")
		this.Redirect(postMd.Link(), 302)
	}
	// update recent/home/category/topics posts cache
	if setting.MemcachedEnabled {
		cache.MemcachedRemove("recent-posts-count")
		cache.MemcachedRemove("recent-posts")
		cache.MemcachedRemove("recent-category")
		cache.MemcachedRemove("home-posts")
		cache.MemcachedRemove("today-topten-posts")
		categoryCountKey := fmt.Sprintf(`category-%s-count`, postMd.Category.Slug)
		cache.MemcachedRemove(categoryCountKey)
		categoryKey := fmt.Sprintf(`category-%s`, postMd.Category.Slug)
		cache.MemcachedRemove(categoryKey)
		topicCountKey := fmt.Sprintf(`topic-%s-count`, postMd.Topic.Slug)
		cache.MemcachedRemove(topicCountKey)
		topicKey := fmt.Sprintf(`topic-%s`, postMd.Topic.Slug)
		cache.MemcachedRemove(topicKey)
	}
	if setting.RedisEnabled {
		cache.RedisRemove("recent-posts-count")
		cache.RedisRemove("recent-posts")
		cache.RedisRemove("recent-category")
		cache.RedisRemove("home-posts")
		cache.RedisRemove("today-topten-posts")
		categoryCountKey := fmt.Sprintf(`category-%s-count`, postMd.Category.Slug)
		cache.RedisRemove(categoryCountKey)
		categoryKey := fmt.Sprintf(`category-%s`, postMd.Category.Slug)
		cache.RedisRemove(categoryKey)
		topicCountKey := fmt.Sprintf(`topic-%s-count`, postMd.Topic.Slug)
		cache.RedisRemove(topicCountKey)
		topicKey := fmt.Sprintf(`topic-%s`, postMd.Topic.Slug)
		cache.RedisRemove(topicKey)
	}
}

func (this *PostRouter) Append() {
	this.TplNames = "post/append.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	postMd.Content = ""
	postMd.ContentCache = ""
	form := post.PostForm{}
	form.SetFromPost(&postMd)
	this.SetFormSets(&form)
}

func (this *PostRouter) AppendSubmit() {
	this.TplNames = "post/append.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	form := post.PostForm{}
	form.SetFromPost(&postMd)
	if !this.ValidAppendFormSets(&form) {
		return
	}

	if len(postMd.Content) == 0 {
		return
	}
	var appendPostMd models.AppendPost
	appendPostMd.Message = form.Content
	appendPostMd.Post = &postMd

	if err := form.AppendPost(&appendPostMd, &this.User); err == nil {
		this.JsStorage("deleteKey", "post/append")
		this.Redirect(postMd.Link(), 302)
	}
}
