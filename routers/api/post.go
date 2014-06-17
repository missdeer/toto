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

package api

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/missdeer/KellyBackend/modules/models"
)

func (this *ApiRouter) PostToggle() {
    result := map[string]interface{}{
        "success": false,
    }

    defer func() {
        this.Data["json"] = result
        this.ServeJson()
    }()

    if !this.IsAjax() {
        return
    }

    action := this.GetString("action")

    if this.IsLogin {

        switch action {
        case "toggle-best":
            id, _ := this.GetInt("post")
            if id > 0 {
                o := orm.NewOrm()
                p := models.Post{ Id: int(id) }
                o.Read(&p);

                p.IsBest = !p.IsBest
                if _, err := o.Update(&p); err != nil {
                    beego.Error("PostCounterAdd ", err)
                } else {
                    result["success"] = true
                }
            }
        case "toggle-top":
            id, _ := this.GetInt("post")
            if id > 0 {
                o := orm.NewOrm()
                p := models.Post{ Id: int(id) }
                o.Read(&p);

                p.IsTop = !p.IsTop
                if _, err := o.Update(&p); err != nil {
                    beego.Error("PostCounterAdd ", err)
                } else {
                    result["success"] = true
                }
            }
        }
    }
}
