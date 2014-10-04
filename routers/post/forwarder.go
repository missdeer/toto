package post

import (
	"fmt"
	"github.com/missdeer/KellyBackend/routers/base"
)

type ForwarderRouter struct {
	base.BaseRouter
}

func (this *ForwarderRouter) TaobaoItem() {
	id := this.GetString("id")
	this.Redirect(fmt.Sprintf("http://item.taobao.com/item.htm?id=%s", id), 302)
}
