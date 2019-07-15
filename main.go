package main

import (
	"fmt"
	"time"
	"wxapp/config"
	"wxapp/pkg/template"
	"wxapp/web/route"

	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/boltdb"
	"os"
	"wxapp/models"
)

func newLogFile() *os.File {
	filename := fmt.Sprintf("logs/%s.txt", time.Now().Format("2006-01-02"))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

func main() {
	// 初始化配置
	config.InitConfig()

	f := newLogFile()
	defer f.Close()
	golog.AddOutput(f)

	app := iris.New()
	// 注册模板引擎
	tmpl := iris.HTML("./web/template", ".html")
	tmpl.Reload(config.ProdMode == false)
	template.BindTmpFunc(tmpl)
	app.RegisterView(tmpl)

	app.Logger().SetLevel(config.LogLevel)
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	// 注册session
	session_db, _ := boltdb.New("./sessions/sessions.db", 0666, "users")
	manager := sessions.New(sessions.Config{
		Cookie:  "adbsessionsid",
		Expires: 24 * time.Hour,
	})
	manager.UseDatabase(session_db)

	// 初始化数据库
	models.NewEngine()
	// 初始化缓存
	models.InitCache()

	//退出程序的时候关闭数据库
	iris.RegisterOnInterrupt(func() {
		models.Close()
		session_db.Close()
	})

	// 设置路由
	route.Route(app, manager)
	//app.Any("/debug/pprof/{action:path}", pprof.New())

	// 注册中间价
	app.UseGlobal(logger.New())
	app.UseGlobal(recover.New())
	app.UseGlobal(iris.Gzip)

	addr := iris.Addr(fmt.Sprintf("%v:%v", config.HTTPAddr, config.HTTPPort))
	if config.ProdMode {
		app.Run(addr, iris.WithoutVersionChecker)
	} else {
		app.Run(addr)
	}
}
