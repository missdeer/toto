package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"regexp"
)

func Render(content string) string {
	// select the renderer
	subs := regexp.MustCompile(`^!([A-Za-z0-9]+)!`).FindStringSubmatch(content)
	if len(subs) <= 1 {
		// no custom specified renderer
		return RenderMarkdown(content)
	}
	renderer := subs[1]
	switch renderer {
	case "m":
		beego.Info("User indicates using markdown renderer!")
		return RenderMarkdown(content[3:])
	case "mp":
		beego.Info("User indicates using weixin mp render!")
		matched := regexp.MustCompile(`^!mp!(http|https)\://([a-zA-Z0-9\.\-]+(\:[a-zA-Z0-9\.&amp;%\$\-]+)*@)?((25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])|([a-zA-Z0-9\-]+\.)*[a-zA-Z0-9\-]+\.[a-zA-Z]{2,4})(\:[0-9]+)?(/[^/][a-zA-Z0-9\.\,\?\'\\/\+&amp;%\$#\=~_\-@]*)*$`).MatchString(content)
		if matched {
			raw := content[4:]
			u, err := url.Parse(raw)
			if err != nil {
				beego.Error("parsing URL failed, fallthrough using markdown renderer")
				break
			}
			if u.Host != "mp.weixin.qq.com" {
				beego.Error("not from mp.weixin.qq.com, fallthrough using markdown renderer")
				break
			}
			s := fmt.Sprintf(`以下内容由系统自动提取自<a href="%s" target='_blank'>%s</a>，点击该链接可访问原文，所有权利归原文出处所有。<hr/>`, raw, u.Host)
			return s + renderWeiXinMP(content[4:])
		}
		break
	case "card":
		beego.Info("User indicates using embedly renderer!")
		matched := regexp.MustCompile(`^!card!(http|https)\://([a-zA-Z0-9\.\-]+(\:[a-zA-Z0-9\.&amp;%\$\-]+)*@)?((25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])|([a-zA-Z0-9\-]+\.)*[a-zA-Z0-9\-]+\.[a-zA-Z]{2,4})(\:[0-9]+)?(/[^/][a-zA-Z0-9\.\,\?\'\\/\+&amp;%\$#\=~_\-@]*)*$`).MatchString(content)
		if matched {
			raw := content[6:]
			_, err := url.Parse(raw)
			if err != nil {
				beego.Error("parsing URL failed, fallthrough using markdown renderer")
				break
			}
			//s := fmt.Sprintf(`以下内容由系统自动提取自<a href="%s" target='_blank'>%s</a>，点击该链接可访问原文，所有权利归原文出处所有。<hr/>`, raw, u.Host)
			//return s + fmt.Sprintf(`<iframe src="https://cdn.embedly.com/widgets/card.html?url=%s" allowfullscreen="" frameborder="0" width="100%" ></iframe>`, url.QueryEscape(raw))
			return fmt.Sprintf(`<a class="embedly-card" href="%s"></a><script>!function(a){var b="embedly-platform",c="script";if(!a.getElementById(b)){var d=a.createElement(c);d.id=b,d.src=("https:"===document.location.protocol?"https":"http")+"://cdn.embedly.com/widgets/platform.js";var e=document.getElementsByTagName(c)[0];e.parentNode.insertBefore(d,e)}}(document);</script>`, content[6:])
		}
		break
	case "y":
		beego.Info("User indicates using youku rendereer!")
		matched := regexp.MustCompile(`^!y!(http|https)\://([a-zA-Z0-9\.\-]+(\:[a-zA-Z0-9\.&amp;%\$\-]+)*@)?((25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])|([a-zA-Z0-9\-]+\.)*[a-zA-Z0-9\-]+\.[a-zA-Z]{2,4})(\:[0-9]+)?(/[^/][a-zA-Z0-9\.\,\?\'\\/\+&amp;%\$#\=~_\-@]*)*$`).MatchString(content)
		if matched {
			raw := content[3:]
			u, err := url.Parse(raw)
			if err != nil {
				beego.Error("parsing URL failed, fallthrough using markdown renderer")
				break
			}
			if u.Host != "v.youku.com" {
				beego.Error("not from youku.com, fallthrough using markdown renderer")
				break
			}
			s := fmt.Sprintf(`以下内容由系统自动提取自<a href="%s" target='_blank'>%s</a>，点击该链接可访问原文，所有权利归原文出处所有。<hr/>`, raw, u.Host)
			ids := regexp.MustCompile(`id_([A-Za-z0-9_]+)\.html`).FindStringSubmatch(raw)
			if len(ids) <= 1 {
				beego.Error("can't find the correct id, fallthrough using markdonw renderer!")
				break
			}
			id := ids[1]
			return s + fmt.Sprintf(`<iframe src="http://player.youku.com/embed/%s" allowfullscreen="" frameborder="0" width="100%" height="420"></iframe>`, id)
		}
		break
	case "r", "e", "p":
		beego.Info("User indicates using readability renderer!")
		matched := regexp.MustCompile(`^!(r|e|p)!(http|https)\://([a-zA-Z0-9\.\-]+(\:[a-zA-Z0-9\.&amp;%\$\-]+)*@)?((25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])|([a-zA-Z0-9\-]+\.)*[a-zA-Z0-9\-]+\.[a-zA-Z]{2,4})(\:[0-9]+)?(/[^/][a-zA-Z0-9\.\,\?\'\\/\+&amp;%\$#\=~_\-@]*)*$`).MatchString(content)
		if matched {
			raw := content[3:]
			u, err := url.Parse(raw)
			if err != nil {
				beego.Error("parsing URL failed, fallthrough using markdown renderer")
				break
			}
			s := fmt.Sprintf(`以下内容由系统自动提取自<a href="%s" target='_blank'>%s</a>，点击该链接可访问原文，所有权利归原文出处所有。<hr/>`, raw, u.Host)
			return s + RenderReadability(raw, renderer)
		}
		break
	default:
		break
	}
	return RenderMarkdown(content)
}
