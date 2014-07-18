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

// Package utils implemented some useful functions.

package setting

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/howeyc/fsnotify"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/captcha"
	"github.com/beego/compress"
	"github.com/beego/i18n"
	"github.com/beego/social-auth"
	"github.com/beego/social-auth/apps"
)

const (
	APP_VER = "1.0.0.1"
)

type AdRecord struct {
	XMLName xml.Name `xml:"ad"`
	Url     string   `xml:"url"`
	Img     []string `xml:"img"`
	Title   string   `xml:"description"`
}

type AdRecords struct {
	XMLName xml.Name   `xml:"yiiliads"`
	Records []AdRecord `xml:"ad"`
}

var (
	AppName             string
	AppVer              string
	AppHost             string
	AppUrl              string
	AppLogo             string
	ImgBedUrl           string
	QiniuEnabled        bool
	QiniuUrl            string
	QiniuAppKey         string
	QiniuSecretKey      string
	QiniuBucketName     string
	UpYunEnabled        bool
	UpYunUrl            string
	UpYunUsername       string
	UpYunPassword       string
	UpYunBucketName     string
	EnforceRedirect     bool
	AvatarURL           string
	SecretKey           string
	ORCAVerifyCode      string
	IsProMode           bool
	ActiveCodeLives     int
	ResetPwdCodeLives   int
	CategoryHideOnHome  int
	DateFormat          string
	DateTimeFormat      string
	DateTimeShortFormat string
	TimeZone            string
	RealtimeRenderMD    bool
	ImageSizeSmall      int
	ImageSizeMiddle     int
	ImageLinkAlphabets  []byte
	ImageXSend          bool
	ImageXSendHeader    string
	Langs               []string

	LoginRememberDays int
	LoginMaxRetries   int
	LoginFailedBlocks int

	CookieRememberName string
	CookieUserName     string

	// search
	SearchEnabled bool
	NativeSearch  bool

	// sphinx search setting
	SphinxEnabled bool
	SphinxHost    string
	SphinxIndex   string
	SphinxMaxConn int

	// mail setting
	MailUser     string
	MailFrom     string
	MailHost     string
	MailAuthUser string
	MailAuthPass string

	// ads setting
	Ads AdRecords

	// memcached setting
	MemcachedEnabled bool
	MemcachedConn    string

	// redis setting
	RedisEnabled bool
	RedisConn    string
)

var (
	// OAuth
	GithubClientId       string
	GithubClientSecret   string
	GoogleClientId       string
	GoogleClientSecret   string
	WeiboClientId        string
	WeiboClientSecret    string
	QQClientId           string
	QQClientSecret       string
	TwitterClientId      string
	TwitterClientSecret  string
	FacebookClientId     string
	FacebookClientSecret string
	FanfouClientId       string
	FanfouClientSecret   string
	SohuClientId         string
	SohuClientSecret     string
	NeteaseClientId      string
	NeteaseClientSecret  string
)

const (
	LangEnUS = iota
	LangZhCN
)

var (
	// Social Auth
	GithubAuth *apps.Github
	GoogleAuth *apps.Google
	SocialAuth *social.SocialAuth
)

var (
	Cfg     *goconfig.ConfigFile
	Cache   cache.Cache
	Captcha *captcha.Captcha
)

var (
	GlobalConfPath   = "conf/global/app.ini"
	AppConfPath      = "conf/app.ini"
	CompressConfPath = "conf/compress.json"
	AdsConfPath      = "conf/ads.xml"
)

func LoadAds() {
	fh, err := os.Open(AdsConfPath)
	if err != nil {
		beego.Error(err)
		return
	}
	defer fh.Close()
	adsxml, err := ioutil.ReadAll(fh)
	if err != nil {
		beego.Error(err)
		return
	}

	Ads.Records = []AdRecord{}
	err = xml.Unmarshal(adsxml, &Ads)
	if err != nil {
		beego.Error(err)
		return
	}

	//for _, ad := range Ads.Records {
	//	fmt.Printf("url: %s, title: %s\n", ad.Url, ad.Title)
	//	for _, img := range ad.Img {
	//		fmt.Printf("img: %s\n", img)
	//	}
	//}
}

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and log level.
	Cfg, err = goconfig.LoadConfigFile(GlobalConfPath)

	if Cfg == nil {
		Cfg, err = goconfig.LoadConfigFile(AppConfPath)
		if err != nil {
			fmt.Println("Fail to load configuration file: " + err.Error())
			os.Exit(2)
		}

	} else {
		Cfg.AppendFiles(AppConfPath)
	}

	Cfg.BlockMode = false

	// set time zone of system
	TimeZone = Cfg.MustValue("app", "time_zone", "Asia/Shanghai")
	if _, err := time.LoadLocation(TimeZone); err == nil {
		os.Setenv("TZ", TimeZone)
	} else {
		fmt.Println("Wrong time_zone: " + TimeZone + " " + err.Error())
		os.Exit(2)
	}

	// Trim 4th part.
	AppVer = strings.Join(strings.Split(APP_VER, ".")[:3], ".")

	beego.RunMode = Cfg.MustValue("app", "run_mode")
	beego.HttpPort = Cfg.MustInt("app", "http_port")

	IsProMode = beego.RunMode == "pro"
	if IsProMode {
		beego.SetLevel(beego.LevelInfo)
	}

	// cache system
	Cache, err = cache.NewCache("memory", `{"interval":360}`)

	Captcha = captcha.NewCaptcha("/captcha/", Cache)
	Captcha.FieldIdName = "CaptchaId"
	Captcha.FieldCaptchaName = "Captcha"

	// session settings
	beego.SessionOn = true
	beego.SessionProvider = Cfg.MustValue("session", "session_provider", "file")
	beego.SessionSavePath = Cfg.MustValue("session", "session_path", "sessions")
	beego.SessionName = Cfg.MustValue("session", "session_name", "yiili_sess")
	beego.SessionCookieLifeTime = Cfg.MustInt("session", "session_life_time", 0)
	beego.SessionGCMaxLifetime = Cfg.MustInt64("session", "session_gc_time", 86400)

	beego.EnableXSRF = true
	// xsrf token expire time
	beego.XSRFExpire = 86400 * 365

	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	dataSource := Cfg.MustValue("orm", "data_source", "root:root@/yiili?charset=utf8&loc=Local")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)

	// set default database
	err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	if err != nil {
		beego.Error(err)
	}
	orm.RunCommand()

	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		beego.Error(err)
	}

	reloadConfig()

	if SphinxEnabled {
		// for search config
		SphinxHost = Cfg.MustValue("search", "sphinx_host", "127.0.0.1:9306")
		SphinxMaxConn = Cfg.MustInt("search", "sphinx_max_conn", 5)
		orm.RegisterDriver("sphinx", orm.DR_MySQL)
	}

	if MemcachedEnabled {
		MemcachedConn = Cfg.MustValue("memcached", "conn", "127.0.0.1:11211")
	}

	if RedisEnabled {
		RedisConn = Cfg.MustValue("redis", "conn", "127.0.0.1:6379")
	}

	social.DefaultAppUrl = AppUrl

	// OAuth
	var clientId, secret string

	clientId = Cfg.MustValue("oauth", "github_client_id", "your_client_id")
	secret = Cfg.MustValue("oauth", "github_client_secret", "your_client_secret")
	GithubAuth = apps.NewGithub(clientId, secret)

	clientId = Cfg.MustValue("oauth", "google_client_id", "your_client_id")
	secret = Cfg.MustValue("oauth", "google_client_secret", "your_client_secret")
	GoogleAuth = apps.NewGoogle(clientId, secret)

	err = social.RegisterProvider(GithubAuth)
	if err != nil {
		beego.Error(err)
	}
	err = social.RegisterProvider(GoogleAuth)
	if err != nil {
		beego.Error(err)
	}

	settingLocales()
	settingCompress()

	configWatcher()

	return Cfg
}

func reloadConfig() {
	AppName = Cfg.MustValue("app", "app_name", "")
	beego.AppName = AppName

	AppHost = Cfg.MustValue("app", "app_host", "127.0.0.1:8092")
	AppUrl = Cfg.MustValue("app", "app_url", "http://127.0.0.1:8092/")
	AppLogo = Cfg.MustValue("app", "app_logo", "")
	ImgBedUrl = Cfg.MustValue("app", "imgbed_url", "")
	QiniuEnabled = Cfg.MustBool("app", "qiniu_enabled", true)
	if QiniuEnabled {
		QiniuUrl = Cfg.MustValue("app", "qiniu_url", "")
		QiniuAppKey = Cfg.MustValue("app", "qiniu_appkey", "")
		QiniuSecretKey = Cfg.MustValue("app", "qiniu_secretkey", "")
		QiniuBucketName = Cfg.MustValue("app", "qiniu_bucketname", "")
	}
	UpYunEnabled = Cfg.MustBool("app", "upyun_enabled", false)
	if UpYunEnabled {
		UpYunUrl = Cfg.MustValue("app", "upyun_url", "")
		UpYunUsername = Cfg.MustValue("app", "upyun_username", "")
		UpYunPassword = Cfg.MustValue("app", "upyun_password", "")
		UpYunBucketName = Cfg.MustValue("app", "upyun_bucketname", "")
	}
	AvatarURL = Cfg.MustValue("app", "avatar_url")

	CategoryHideOnHome = Cfg.MustInt("app", "category_hide_on_home", 0)
	EnforceRedirect = Cfg.MustBool("app", "enforce_redirect")

	DateFormat = Cfg.MustValue("app", "date_format")
	DateTimeFormat = Cfg.MustValue("app", "datetime_format")
	DateTimeShortFormat = Cfg.MustValue("app", "datetime_short_format")

	SecretKey = Cfg.MustValue("app", "secret_key")
	if len(SecretKey) == 0 {
		fmt.Println("Please set your secret_key in app.ini file")
	}

	ORCAVerifyCode = Cfg.MustValue("app", "orca_verify_code")

	ActiveCodeLives = Cfg.MustInt("app", "acitve_code_live_minutes", 180)
	ResetPwdCodeLives = Cfg.MustInt("app", "resetpwd_code_live_minutes", 180)

	LoginRememberDays = Cfg.MustInt("app", "login_remember_days", 7)
	LoginMaxRetries = Cfg.MustInt("app", "login_max_retries", 5)
	LoginFailedBlocks = Cfg.MustInt("app", "login_failed_blocks", 10)

	CookieRememberName = Cfg.MustValue("app", "cookie_remember_name", "yiili_magic")
	CookieUserName = Cfg.MustValue("app", "cookie_user_name", "yiili_powerful")

	RealtimeRenderMD = Cfg.MustBool("app", "realtime_render_markdown")

	ImageSizeSmall = Cfg.MustInt("image", "image_size_small")
	ImageSizeMiddle = Cfg.MustInt("image", "image_size_middle")

	if ImageSizeSmall <= 0 {
		ImageSizeSmall = 300
	}

	if ImageSizeMiddle <= ImageSizeSmall {
		ImageSizeMiddle = ImageSizeSmall + 400
	}

	str := Cfg.MustValue("image", "image_link_alphabets")
	if len(str) == 0 {
		str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	ImageLinkAlphabets = []byte(str)

	ImageXSend = Cfg.MustBool("image", "image_xsend", false)
	ImageXSendHeader = Cfg.MustValue("image", "image_xsend_header", "X-Accel-Redirect")

	MailUser = Cfg.MustValue("mailer", "mail_name", "Yiili Community")
	MailFrom = Cfg.MustValue("mailer", "mail_from", "noreply@yii.li")

	// set mailer connect args
	MailHost = Cfg.MustValue("mailer", "mail_host", "127.0.0.1:25")
	MailAuthUser = Cfg.MustValue("mailer", "mail_user", "example@example.com")
	MailAuthPass = Cfg.MustValue("mailer", "mail_pass", "******")

	orm.Debug = Cfg.MustBool("orm", "debug_log")

	// search setting
	SphinxIndex = Cfg.MustValue("search", "sphinx_index", "yiili, yiili_delta")

	SearchEnabled = Cfg.MustBool("search", "enabled")
	SphinxEnabled = Cfg.MustBool("search", "sphinx_enabled")
	NativeSearch = Cfg.MustBool("search", "native_search")
	if !SearchEnabled {
		SphinxEnabled = false
		NativeSearch = false
	}

	// memcached, redis
	MemcachedEnabled = Cfg.MustBool("memcached", "enabled")
	RedisEnabled = Cfg.MustBool("redis", "enabled")

	// OAuth
	GithubClientId = Cfg.MustValue("oauth", "github_client_id", "your_client_id")
	GithubClientSecret = Cfg.MustValue("oauth", "github_client_secret", "your_client_secret")
	GoogleClientId = Cfg.MustValue("oauth", "google_client_id", "your_client_id")
	GoogleClientSecret = Cfg.MustValue("oauth", "google_client_secret", "your_client_secret")
	WeiboClientId = Cfg.MustValue("oauth", "weibo_client_id", "your_client_id")
	WeiboClientSecret = Cfg.MustValue("oauth", "weibo_client_secret", "your_client_secret")
	QQClientId = Cfg.MustValue("oauth", "qq_client_id", "your_client_id")
	QQClientSecret = Cfg.MustValue("oauth", "qq_client_secret", "your_client_secret")
	TwitterClientId = Cfg.MustValue("oauth", "twitter_client_id", "your_client_id")
	TwitterClientSecret = Cfg.MustValue("oauth", "twitter_client_secret", "your_client_secret")
	FacebookClientId = Cfg.MustValue("oauth", "facebook_client_id", "your_client_id")
	FacebookClientSecret = Cfg.MustValue("oauth", "facebook_client_secret", "your_client_secret")
	FanfouClientId = Cfg.MustValue("oauth", "fanfou_client_id", "your_client_id")
	FanfouClientSecret = Cfg.MustValue("oauth", "fanfou_client_secret", "your_client_secret")
	SohuClientId = Cfg.MustValue("oauth", "sohu_client_id", "your_client_id")
	SohuClientSecret = Cfg.MustValue("oauth", "sohu_client_secret", "your_client_secret")
	NeteaseClientId = Cfg.MustValue("oauth", "netease_client_id", "your_client_id")
	NeteaseClientSecret = Cfg.MustValue("oauth", "netease_client_secret", "your_client_secret")
}

func settingLocales() {
	// load locales with locale_LANG.ini files
	langs := "zh-CN|en-US"
	for _, lang := range strings.Split(langs, "|") {
		lang = strings.TrimSpace(lang)
		files := []string{"conf/" + "locale_" + lang + ".ini"}
		if fh, err := os.Open(files[0]); err == nil {
			fh.Close()
		} else {
			files = nil
		}
		if err := i18n.SetMessage(lang, "conf/global/"+"locale_"+lang+".ini", files...); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			os.Exit(2)
		}
	}
	Langs = i18n.ListLangs()
}

func settingCompress() {
	setting, err := compress.LoadJsonConf(CompressConfPath, IsProMode, ImgBedUrl)
	if err != nil {
		beego.Error(err)
		return
	}

	setting.RunCommand()

	if IsProMode {
		setting.RunCompress(true, false, true)
	}

	beego.AddFuncMap("compress_js", setting.Js.CompressJs)
	beego.AddFuncMap("compress_css", setting.Css.CompressCss)
}

var eventTime = make(map[string]int64)

func configWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("Failed start app watcher: " + err.Error())
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				switch filepath.Ext(event.Name) {
				case ".ini":
					if checkEventTime(event.Name) {
						continue
					}
					beego.Info(event)

					if err := Cfg.Reload(); err != nil {
						beego.Error("Conf Reload: ", err)
					}

					if err := i18n.ReloadLangs(); err != nil {
						beego.Error("Conf Reload: ", err)
					}

					reloadConfig()
					beego.Info("Config Reloaded")

				case ".xml":
					if checkEventTime(event.Name) {
						continue
					}

					LoadAds()
					beego.Info("Ads config reloaded")
				case ".json":
					if checkEventTime(event.Name) {
						continue
					}
					if event.Name == CompressConfPath {
						settingCompress()
						beego.Info("Beego Compress Reloaded")
					}
				}
			}
		}
	}()

	if err := watcher.WatchFlags("conf", fsnotify.FSN_MODIFY); err != nil {
		beego.Error(err)
	}

	if err := watcher.WatchFlags("conf/global", fsnotify.FSN_MODIFY); err != nil {
		beego.Error(err)
	}
}

// checkEventTime returns true if FileModTime does not change.
func checkEventTime(name string) bool {
	mt := getFileModTime(name)
	if eventTime[name] == mt {
		return true
	}

	eventTime[name] = mt
	return false
}

// getFileModTime retuens unix timestamp of `os.File.ModTime` by given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		beego.Error("Fail to open file[ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		beego.Error("Fail to get file information[ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

func IsMatchHost(uri string) bool {
	if len(uri) == 0 {
		return false
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	if u.Host != AppHost {
		return false
	}

	return true
}
