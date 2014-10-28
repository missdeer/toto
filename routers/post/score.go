package post

import (
	"github.com/missdeer/toto/routers/base"
)

// HomeRouter serves home page.
type ScoreRouter struct {
	base.BaseRouter
}

func (this *ScoreRouter) Home() {
	this.Data["IsHome"] = false
	this.TplNames = "post/score.html"
}
