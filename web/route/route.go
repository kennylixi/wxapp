package route

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/hero"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"time"
	"wxapp/config"
	"wxapp/web/controller"
	"wxapp/web/controller/admin"
	"wxapp/web/controller/front"
)

var (
	smanager *sessions.Sessions
	crs      context.Handler
)

const (
	expiresEvery = 5 * 60 * time.Second
)

func Route(app *iris.Application, sess *sessions.Sessions) {
	smanager = sess
	//跨域访问
	crs = cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hoss.
		AllowCredentials: true,
	})

	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)

	registFrontRoute(app)

	// 注册管理后台路由
	registAdminRoute(app)
}

func registFrontRoute(app *iris.Application) {
	frontParty := app.Party("/", crs).AllowMethods(iris.MethodOptions)
	frontParty.StaticWeb("/", "./static")
	if !config.ProdMode {
		frontParty.StaticWeb("/static", "./web/static")
	}
	//直接判断一级路由中间是否存在点，如果存在直接当做文件来处理
	//frontParty.Use(func(ctx context.Context) {
	//	path := ctx.Path()
	//
	//	if strings.Index(path,".")!=-1 {
	//		ctx.ServeFile(fmt.Sprintf("./static/%s",path),false)
	//		return
	//	}
	//	ctx.Next()
	//})
	//if config.ProdMode {
	//	//添加403缓存
	//	frontParty.Use(func(ctx context.Context) {
	//		path := ctx.Path()
	//
	//		if path == "/apply" || path == "/search" || path == "/captcha" {
	//			//filter一些不需要缓存的
	//		} else {
	//			now := time.Now()
	//			if modified, err := ctx.CheckIfModifiedSince(now.Add(-expiresEvery)); !modified && err == nil {
	//				ctx.WriteNotModified()
	//				return
	//			}
	//
	//			ctx.SetLastModified(now)
	//		}
	//		ctx.Next()
	//	})
	//}

	mvc.Configure(frontParty.Party("/"), func(app *mvc.Application) {
		app.Handle(new(front.HomeController))
	})
	mvc.Configure(frontParty.Party("/upload"), func(app *mvc.Application) {
		app.Handle(new(controller.UploadController))
	})
}

func registAdminRoute(app *iris.Application) {
	adminParty := app.Party("admin.", crs).AllowMethods(iris.MethodOptions)
	if !config.ProdMode {
		adminParty.StaticWeb("/static", "./web/static")
	}
	adminParty.Use(func(context iris.Context) {
		path := context.Path()

		if path == "/login" {
			//filter一些不需要登录的
		} else {
			//没有登录跳转到登录
			userId := smanager.Start(context).GetInt64Default(controller.SESSION_USER_ID_KEY, 0)
			if userId <= 0 {
				context.Redirect("/login")
			}
		}

		context.Next()
	})
	hero.Register(smanager.Start)
	mvc.Configure(adminParty.Party("/"), func(app *mvc.Application) {
		app.Handle(new(admin.HomeController))
	})
	mvc.Configure(adminParty.Party("/user"), func(app *mvc.Application) {
		app.Handle(new(admin.UserController))
	})
	mvc.Configure(adminParty.Party("/article"), func(app *mvc.Application) {
		app.Handle(new(admin.ArticleController))
	})
	mvc.Configure(adminParty.Party("/category"), func(app *mvc.Application) {
		app.Handle(new(admin.CategoryController))
	})
	mvc.Configure(adminParty.Party("/city"), func(app *mvc.Application) {
		app.Handle(new(admin.CityController))
	})
	mvc.Configure(adminParty.Party("/scorelog"), func(app *mvc.Application) {
		app.Handle(new(admin.ScorelogController))
	})
	mvc.Configure(adminParty.Party("/special"), func(app *mvc.Application) {
		app.Handle(new(admin.SpecialController))
	})
	mvc.Configure(adminParty.Party("/upload"), func(app *mvc.Application) {
		app.Handle(new(controller.UploadController))
	})
	mvc.Configure(adminParty.Party("/link"), func(app *mvc.Application) {
		app.Handle(new(admin.LinkController))
	})
	mvc.Configure(adminParty.Party("/apply"), func(app *mvc.Application) {
		app.Handle(new(admin.ApplyController))
	})
	mvc.Configure(adminParty.Party("/recommend"), func(app *mvc.Application) {
		app.Handle(new(admin.RecommendController))
	})
	mvc.Configure(adminParty.Party("/evaluate"), func(app *mvc.Application) {
		app.Handle(new(admin.EvaluateController))
	})
	mvc.Configure(adminParty.Party("/mp"), func(app *mvc.Application) {
		app.Handle(new(admin.MpController))
	})
}

func notFoundHandler(ctx iris.Context) {
	ctx.View("front/common/404.html")
}
