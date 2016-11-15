// Copyright 2014 missdeer
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

package post

import (
	"github.com/missdeer/toto/routers/base"
)

// HomeRouter serves home page.
type AuxiliaryRouter struct {
	base.BaseRouter
}

func (this *AuxiliaryRouter) Err401() {
	this.TplName = "post/error/401.html"
}

func (this *AuxiliaryRouter) Err403() {
	this.TplName = "post/error/403.html"
}

func (this *AuxiliaryRouter) Err404() {
	this.TplName = "post/error/404.html"
}

func (this *AuxiliaryRouter) Err500() {
	this.TplName = "post/error/500.html"
}

func (this *AuxiliaryRouter) Err503() {
	this.TplName = "post/error/503.html"
}
