package config

import (
	"fmt"
	"github.com/kataras/golog"
	"gopkg.in/ini.v1"
	"wxapp/pkg/bindata"
)

var (
	AppName   string
	AppUrl    string
	StaticUrl string
	Domain    string
	HTTPAddr  string
	HTTPPort  string

	//七牛
	QiniuBucket    string
	QAccessKey     string
	QSecretKey     string
	QPipeline      string
	QPersistentOps string
	QDomain        string

	SeoBaiduPushToken string//百度push Token

	Cfg      *ini.File
	ProdMode bool

	LogLevel string
)

func init() {
}

func InitConfig() {
	//iniPath, err := filepath.Abs("./config/app.ini")
	var err error
	Cfg, err = ini.Load(bindata.MustAsset("config/app.ini"))
	if err != nil {
		golog.Fatalf("Fail to parse 'config/app.ini': %v", err)
	}
	Cfg.NameMapper = ini.AllCapsUnderscore

	AppName = Cfg.Section("").Key("APP_NAME").MustString("Apd")
	ProdMode = Cfg.Section("").Key("RUN_MODE").String() == "prod"
	StaticUrl = Cfg.Section("").Key("STATIC_DOMAIN").MustString("")

	sec := Cfg.Section("server")
	Domain = sec.Key("DOMAIN").MustString("localhost")
	HTTPAddr = sec.Key("HTTP_ADDR").MustString("0.0.0.0")
	HTTPPort = sec.Key("HTTP_PORT").MustString("80")
	AppUrl = fmt.Sprintf("http://%s", Domain)
	//日志
	logSec := Cfg.Section("log")
	LogLevel = logSec.Key("LEVEL").MustString("debug")
	//七牛
	qiniu := Cfg.Section("qiniu")
	QiniuBucket = qiniu.Key("BUCKET").MustString("")
	QAccessKey = qiniu.Key("ACCESS_KEY").MustString("")
	QSecretKey = qiniu.Key("SECRET_KEY").MustString("")
	QPipeline = qiniu.Key("PIPELINE").MustString("")
	QPersistentOps = qiniu.Key("PERSISTENT_OPS").MustString("")
	QDomain = qiniu.Key("DOMAIN").MustString("")
	//seo
	seo := Cfg.Section("seo")
	SeoBaiduPushToken = seo.Key("BAIDU_PUSH_TOKEN").MustString("")
}
