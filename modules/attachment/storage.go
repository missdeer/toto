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
	"fmt"
	"github.com/missdeer/KellyBackend/setting"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"

	"github.com/astaxie/beego"

	"github.com/missdeer/KellyBackend/modules/models"
	"github.com/missdeer/KellyBackend/modules/utils"
	"github.com/missdeer/KellyBackend/upyun"
	. "github.com/qiniu/api/conf"
	qiniuio "github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
)

func uploadToQiniu(imageLocalPath string, imageRemotePath string, result chan error) {
	if setting.QiniuEnabled {
		beego.Info("start uploading to qiniu ", imageLocalPath)
		putPolicy := rs.PutPolicy{}
		putPolicy.Scope = setting.QiniuBucketName
		uptoken := putPolicy.Token(nil)

		var ret qiniuio.PutRet
		var extra = &qiniuio.PutExtra{}

		// get encoded file name as the key
		err := qiniuio.PutFile(nil, &ret, uptoken, imageRemotePath, imageLocalPath, extra)
		if err != nil {
			beego.Error("putting file without key to Qiniu failed: ", err)
		}
		result <- err
	} else {
		result <- nil
	}
	beego.Info("uploaded to qiniu ", imageRemotePath)
}

func uploadToUpyun(imageLocalPath string, imageRemotePath string, result chan error) {
	if setting.UpYunEnabled {
		beego.Info("start uploading to upyun ", imageLocalPath)
		var upyunio *upyun.UpYun
		upyunio = upyun.NewUpYun(setting.UpYunBucketName, setting.UpYunUsername, setting.UpYunPassword)
		f, err := os.OpenFile(imageLocalPath, os.O_RDONLY, 0644)
		if err != nil {
			beego.Error("opening local saved path failed: ", err)
			return
		}
		defer f.Close()
		err = upyunio.WriteFile(imageRemotePath, f, true)
		if err != nil {
			beego.Error("writing file to UpYun failed: ", err)
		}
		result <- err
	} else {
		result <- nil
	}
	beego.Info("uploaded to upyun ", imageRemotePath)
}

func SaveImage(m *models.Image, r io.ReadSeeker, mime string, filename string, created time.Time) error {
	var ext string

	// test image mime type
	switch mime {
	case "image/jpeg":
		ext = ".jpg"

	case "image/png":
		ext = ".png"

	case "image/gif":
		ext = ".gif"

	default:
		ext = filepath.Ext(filename)
		switch ext {
		case ".jpg", ".png", ".gif":
		default:
			return fmt.Errorf("unsupport image format `%s`", filename)
		}
	}

	// decode image
	var img image.Image
	var err error
	switch ext {
	case ".jpg":
		m.Ext = 1
		img, err = jpeg.Decode(r)
	case ".png":
		m.Ext = 2
		img, err = png.Decode(r)
	case ".gif":
		m.Ext = 3
		img, err = gif.Decode(r)
	}

	if err != nil {
		return err
	}

	m.Width = img.Bounds().Dx()
	m.Height = img.Bounds().Dy()
	m.Created = created

	if err := m.Insert(); err != nil || m.Id <= 0 {
		return err
	}

	path := GenImagePath(m)
	os.MkdirAll(path, 0755)

	fullPath := GenImageFilePath(m, 0)
	var smallPath, middlePath string
	if _, err := r.Seek(0, 0); err != nil {
		return err
	}

	var file *os.File
	if f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		file = f
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		os.RemoveAll(fullPath)
		return err
	}

	var key = "upload" + m.LinkFull()

	qiniuFull := make(chan error)
	qiniuSmall := make(chan error)
	qiniuMiddle := make(chan error)
	upyunFull := make(chan error)
	upyunSmall := make(chan error)
	upyunMiddle := make(chan error)

	ACCESS_KEY = setting.QiniuAppKey
	SECRET_KEY = setting.QiniuSecretKey
	go uploadToQiniu(fullPath, key, qiniuFull)
	go uploadToUpyun(fullPath, "/"+key, upyunFull)
	var result error = nil

	if ext != ".gif" {
		if m.Width > setting.ImageSizeSmall {
			if err := ImageResize(m, img, setting.ImageSizeSmall); err != nil {
				<-qiniuFull
				<-upyunFull
				os.RemoveAll(fullPath)
				return err
			}
			smallPath = GenImageFilePath(m, setting.ImageSizeSmall)
			key = "upload" + m.LinkSmall()
			go uploadToQiniu(smallPath, key, qiniuSmall)
			go uploadToUpyun(smallPath, "/"+key, upyunSmall)
		}

		if m.Width > setting.ImageSizeMiddle {
			if err := ImageResize(m, img, setting.ImageSizeMiddle); err != nil {
				<-qiniuFull
				<-upyunFull
				<-qiniuSmall
				<-upyunSmall
				os.RemoveAll(fullPath)
				os.RemoveAll(smallPath)
				return err
			}
			middlePath = GenImageFilePath(m, setting.ImageSizeMiddle)
			key = "upload" + m.LinkMiddle()
			go uploadToQiniu(middlePath, key, qiniuMiddle)
			go uploadToUpyun(middlePath, "/"+key, upyunMiddle)
		}
	}
	qiniuErr, upyunErr := <-qiniuFull, <-upyunFull
	if qiniuErr != nil {
		result = qiniuErr
	}
	if upyunErr != nil {
		result = upyunErr
	}
	os.RemoveAll(fullPath)
	if len(smallPath) > 0 {
		qiniuErr, upyunErr = <-qiniuSmall, <-upyunSmall
		if qiniuErr != nil {
			result = qiniuErr
		}
		if upyunErr != nil {
			result = upyunErr
		}
		os.RemoveAll(smallPath)
	}

	if len(middlePath) > 0 {
		qiniuErr, upyunErr = <-qiniuMiddle, <-upyunMiddle
		if qiniuErr != nil {
			result = qiniuErr
		}
		if upyunErr != nil {
			result = upyunErr
		}
		os.RemoveAll(middlePath)
	}

	return result
}

func ImageResize(img *models.Image, im image.Image, width int) error {
	savePath := GenImageFilePath(img, width)
	im = resize.Resize(uint(width), 0, im, resize.Bilinear)

	var file *os.File
	if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		file = f
	}
	defer file.Close()

	var err error
	switch img.Ext {
	case 1:
		err = jpeg.Encode(file, im, &jpeg.Options{90})
	case 2:
		err = png.Encode(file, im)
	default:
		return fmt.Errorf("<ImageResize> unsupport image format")
	}

	return err
}

func GenImagePath(img *models.Image) string {
	return "upload/img/" + beego.Date(img.Created, "y/m/d/s/") + utils.ToStr(img.Id) + "/"
}

func GenImageFilePath(img *models.Image, width int) string {
	var size string
	if width == 0 {
		size = "full"
	} else {
		size = utils.ToStr(width)
	}
	return GenImagePath(img) + size + img.GetExt()
}
