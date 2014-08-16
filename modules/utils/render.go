package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"regexp"
)

func Render(content string) string {
	// markdown or readability?
	matched, err := regexp.MatchString(`^(http|https)\://([a-zA-Z0-9\.\-]+(\:[a-zA-Z0-9\.&amp;%\$\-]+)*@)?((25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])|([a-zA-Z0-9\-]+\.)*[a-zA-Z0-9\-]+\.[a-zA-Z]{2,4})(\:[0-9]+)?(/[^/][a-zA-Z0-9\.\,\?\'\\/\+&amp;%\$#\=~_\-@]*)*$`,
		content)
	//matched, err := regexp.MatchString(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`, content)
	if err == nil && matched {
		u, err := url.Parse(content)
		if err != nil {
			beego.Error("parsing URL failed")
		}
		s := fmt.Sprintf(`以下内容由系统自动提取自<a href="%s" target='_blank'>%s</a>，点击该链接可访问原文，所有权利归原文出处所有。<hr/>`, content, u.Host)
		return s + RenderReadability(content)
	} else {
		return RenderMarkdown(content)
	}
}
