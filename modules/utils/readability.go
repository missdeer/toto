package utils

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/missdeer/toto/setting"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
)

func renderPythonReadability(content string) string {
	cmd := exec.Command("python", "-m", "readability.readability", "-u", content)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		beego.Error("creating stdout pipe failed: ", err)
		return ""
	}
	if err := cmd.Start(); err != nil {
		beego.Error("starting command failed: ", err)
		return ""
	}
	raw, err := ioutil.ReadAll(stdout)
	if err != nil {
		beego.Error("reading stdout from pipe failed: ", err)
		return ""
	}
	if err := cmd.Wait(); err != nil {
		beego.Error("waiting from command failed:", err)
		return ""
	}
	return string(raw)
}

func renderReadability(content string) string {
	appkey := setting.ReadabilityAppKey
	req := fmt.Sprintf("https://www.readability.com/api/content/v1/parser?url=%s&token=%s", url.QueryEscape(content), appkey)
	resp, err := http.Get(req)
	if err != nil {
		beego.Error("read response error: ", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	json.Unmarshal(body, &res)
	if str, ok := res["content"].(string); ok {
		return str
	} else {
		beego.Error("converting content to string failed")
		return ""
	}
}

func renderEmbedlyExtract(content string) string {
	appkey := setting.EmbedlyAppKey
	req := fmt.Sprintf("http://api.embed.ly/1/extract?key=%s&url=%s", appkey, url.QueryEscape(content))
	resp, err := http.Get(req)
	if err != nil {
		beego.Error("read response error: ", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		beego.Error("json unmarshalling failed: ", err)
		return ""
	}
	fmt.Printf("%v\n", res)
	if str, ok := res["content"].(string); ok {
		return str
	} else {
		beego.Error("converting content to string failed")
		return ""
	}
}

func RenderReadability(content string, interpreter string) string {
	if len(interpreter) == 0 {
		interpreter = setting.ReadabilityBackend
	}
	switch interpreter {
	case "readability", "r":
		beego.Info("readability backend")
		return renderReadability(content)
	case "embedly", "e":
		beego.Info("embedly extract backend")
		return renderEmbedlyExtract(content)
	default:
		beego.Info("default backend")
		return renderPythonReadability(content)
	}
	beego.Info("default backend")
	return renderPythonReadability(content)
}
