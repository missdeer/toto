package post

import (
	"github.com/missdeer/toto/routers/base"
)

// HomeRouter serves home page.
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
