package route

import (
	"github.com/kataras/iris"
	"52jolynn.com/core"
	"52jolynn.com/api"
)

func RegisterRoutes(ctx core.Context, app *iris.Application) {
	mapi := api.NewMapi(ctx)
	app.Get("/index", func(ctx iris.Context) {
		ctx.Writef("Hello assassin")
	})
	v1 := app.Party("/v1")

	{
		v1.Post("/club", func(ctx iris.Context) {
			mapi.CreateClub(ctx.PostValue("name"), nil, nil, nil)
		})
	}
}
