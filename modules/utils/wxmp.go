package utils

import (
	"bytes"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

func renderWeiXinMP(content string) string {
	resp, err := http.Get(content)
	if err != nil {
		beego.Error("read response error: ", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// extract the content body
	startPos := bytes.Index(body, []byte(`<div id="page-content">`))
	endPos := bytes.Index(body[startPos:], []byte(`<script type="text/javascript">`))

	return string(body[startPos:startPos+endPos]) + `</div></div>`
}
