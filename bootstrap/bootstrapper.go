package bootstrap

import (
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/sessions"
	"lottery/conf"
	"lottery/cron"
	"time"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName string
	AppOwner string
	AppSpawnDate time.Time

	Sessions *sessions.Sessions
}

func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Application : iris.New(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}


func (b *Bootstrapper) SetupViews(viewDir string) {
	htmlEngine := iris.HTML(viewDir, ".html").Layout("shared/layout.html")
	htmlEngine.Reload(true)

	htmlEngine.AddFunc("FromUnixtimeShort", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeformShort)
	})

	htmlEngine.AddFunc("FromUnixtime", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeform)
	})

	b.RegisterView(htmlEngine)
}

func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + b.AppName,
		Expires:   expires,
		Encoding : securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app" : b.AppName,
			"status" : ctx.GetStatusCode(),
			"message" : ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		ctx.View("shared/error.html")
	})
}

func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

func (b *Bootstrapper) setupCron() {
	if conf.RunningCrontabService {
		cron.ConfigureAppOneCron()
	}
	cron.ConfigureAppAllCron()
}

const (
	StaticAssets = "./web/public/"
	Favicon = "favicon.ico"
)

func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews("./web/views")
	b.SetupSessions(24*time.Hour, []byte("the-big-and-secret-fash-key-here"), []byte("lot-secret-of-characters-big-too"))
	b.SetupErrorHandlers()

	b.Favicon(StaticAssets + Favicon)
	b.StaticWeb("/public", StaticAssets)
	b.StaticWeb("/web/public", StaticAssets)

	b.setupCron()

	b.Use(recover.New())
	b.Use(logger.New())

	return b
}


func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}

