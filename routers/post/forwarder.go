package post

import (
	"fmt"
	"github.com/missdeer/toto/routers/base"
	"github.com/missdeer/toto/setting"
)

type ForwarderRouter struct {
	base.BaseRouter
}

func (this *ForwarderRouter) TaobaoItem() {
	id := this.GetString(":id")
	url := fmt.Sprintf("http://item.taobao.com/item.htm?id=%s", id)
	this.Redirect(url, 302)
}

func (this *ForwarderRouter) Favicon() {
	url := fmt.Sprintf("http://%s/static/img/favicon.png", setting.AppHost)
	this.Redirect(url, 302)
}
