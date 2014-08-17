package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"regexp"
)

func Render(content string) string {
	// select the renderer
	renderer := content[:3]
	switch renderer {
	case "!m!":
		beego.Info("User indicates using markdown renderer!")
		return RenderMarkdown(content[3:])
	case "!i!":
		beego.Info("User indicates using image renderer!")
		raw := content[3:]
		u, err := url.Parse(raw)
		if err != nil {
			beego.Error("parsing URL failed, fallthrough using markdown renderer")
			break
		}
		s := fmt.Sprintf(`以下内容由系统自动提取自<a href="%s" target='_blank'>%s</a>，点击该链接可访问原文，所有权利归原文出处所有。<hr/>`, raw, u.Host)
		return s + fmt.Sprintf(`<img src="%s">`, raw)
	case "!r!":
		beego.Info("User indicates using readability renderer!")
		matched, err := regexp.MatchString(`^!r!(http|https)\://([a-zA-Z0-9\.\-]+(\:[a-zA-Z0-9\.&amp;%\$\-]+)*@)?((25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])|([a-zA-Z0-9\-]+\.)*[a-zA-Z0-9\-]+\.[a-zA-Z]{2,4})(\:[0-9]+)?(/[^/][a-zA-Z0-9\.\,\?\'\\/\+&amp;%\$#\=~_\-@]*)*$`,
			content)
		if err == nil && matched {
			raw := content[3:]
			u, err := url.Parse(raw)
			if err != nil {
				beego.Error("parsing URL failed, fallthrough using markdown renderer")
				break
			}
			s := fmt.Sprintf(`以下内容由系统自动提取自<a href="%s" target='_blank'>%s</a>，点击该链接可访问原文，所有权利归原文出处所有。<hr/>`, raw, u.Host)
			return s + RenderReadability(raw)
		}
		fallthrough
	default:
		break
	}
	return RenderMarkdown(content)
}
