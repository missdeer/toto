// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package attachment

import (
	"github.com/missdeer/toto/setting"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"

	"github.com/missdeer/toto/modules/attachment"
	"github.com/missdeer/toto/modules/models"
	"github.com/missdeer/toto/routers/base"
)

type UploadRouter struct {
	base.BaseRouter
}

func (this *UploadRouter) Post() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = &result
		this.ServeJSON()
	}()

	// check permition
	if !this.User.IsActive {
		return
	}

	// get file object
	file, handler, err := this.Ctx.Request.FormFile("image")
	if err != nil {
		return
	}
	defer file.Close()

	t := time.Now()

	image := models.Image{}
	image.User = &this.User

	// get mime type
	mime := handler.Header.Get("Content-Type")

	// save and resize image
	if err := attachment.SaveImage(&image, file, mime, handler.Filename, t); err != nil {
		beego.Error(err)
		return
	}

	result["link"] = image.LinkMiddle()
	result["success"] = true

}

func ImageFilter(ctx *context.Context) {
	if setting.QiniuEnabled && setting.UpYunEnabled {
		if time.Now().UnixNano()%2 == 0 {
			url := setting.QiniuUrl + "upload" + ctx.Request.RequestURI
			http.Redirect(ctx.ResponseWriter, ctx.Request, url, 302)
		} else {
			url := setting.UpYunUrl + "upload" + ctx.Request.RequestURI
			http.Redirect(ctx.ResponseWriter, ctx.Request, url, 302)
		}
		return
	}

	if setting.QiniuEnabled {
		url := setting.QiniuUrl + "upload" + ctx.Request.RequestURI
		http.Redirect(ctx.ResponseWriter, ctx.Request, url, 302)
		return
	}

	if setting.UpYunEnabled {
		url := setting.UpYunUrl + "upload" + ctx.Request.RequestURI
		http.Redirect(ctx.ResponseWriter, ctx.Request, url, 302)
	}
}
