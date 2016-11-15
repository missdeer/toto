package post

import (
	"github.com/missdeer/toto/routers/base"
)

// HomeRouter serves home page.
type NewsRouter struct {
	base.BaseRouter
}

func (this *NewsRouter) Home() {
	this.Data["IsHome"] = false
	this.TplName = "post/news.html"
}
