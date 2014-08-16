package utils

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/missdeer/KellyBackend/setting"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func renderPythonReadability(content string) string {
	cmd := exec.Command("python", "-m", "readability.readability", "-u", content)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		beego.Error(err)
	}
	if err := cmd.Start(); err != nil {
		beego.Error(err)
	}
	raw, err := ioutil.ReadAll(stdout)
	if err != nil {
		beego.Error(err)
	}
	if err := cmd.Wait(); err != nil {
		beego.Error(err)
	}
	return string(raw)
}

func renderReadability(content string) string {
	appkey := setting.ReadabilityAppKey
	req := fmt.Sprintf("https://readability.com/api/content/v1/parser?url=%s&token=%s", content, appkey)
	resp, err := http.Get(req)
	if err != nil {
		beego.Error("read response error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	json.Unmarshal(body, &res)
	if str, ok := res["content"].(string); ok {
		return str
	} else {
		return ""
	}
}

func renderEmbedlyExtract(content string) string {
	appkey := setting.EmbedlyAppKey
	req := fmt.Sprintf("http://api.embed.ly/1/extract?key=%s&url=%s", appkey, content)
	resp, err := http.Get(req)
	if err != nil {
		beego.Error("read response error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	json.Unmarshal(body, &res)
	if str, ok := res["content"].(string); ok {
		return str
	} else {
		return ""
	}
}

func RenderReadability(content string) string {
	switch setting.ReadabilityBackend {
	case "readability":
		return renderReadability(content)
	case "embedly":
		return renderEmbedlyExtract(content)
	default:
		return renderPythonReadability(content)
	}
	return renderPythonReadability(content)
}
