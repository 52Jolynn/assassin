package route

import (
	"github.com/kataras/iris"
	"52jolynn.com/core"
	"52jolynn.com/api"
	"time"
	"52jolynn.com/misc"
)

var mapi api.Mapi
var uapi api.Uapi

func RegisterRoutes(ctx core.Context, app *iris.Application) {
	mapi = api.NewMapi(ctx)
	uapi = api.NewUapi(ctx)
	app.Get("/index", func(ctx iris.Context) {
		ctx.Writef("Hello assassin, today is %s", time.Now().Format(misc.StandardTimeFormatPattern))
	})
	v1 := app.Party("/v1")
	//俱乐部
	{
		v1.Get("/club/{id}", getClub)
		v1.Get("/clubs", getClubs)
		v1.Post("/club", createClub)
		v1.Put("/club", updateClub)
	}
	//场地
	{
		v1.Get("/ground/{id}", getGround)
		v1.Get("/grounds", getGrounds)
		v1.Post("/ground", createGround)
		v1.Put("/ground", updateGround)
	}
	//球队
	{
		v1.Get("/team/{id}", getTeam)
		v1.Get("/teams", getTeams)
		v1.Post("/team", createTeam)
		v1.Put("/team", updateTeam)
	}
	//球员
	{
		v1.Get("/player/{id}", getPlayer)
		v1.Get("/players", getPlayers)
		v1.Post("/player", createPlayer)
		v1.Put("/player", updatePlayer)
	}
	//优惠券
	{
		v1.Get("/coupon/{id}", getCoupon)
		v1.Get("/coupons", getCoupons)
		v1.Post("/coupon", createCoupon)
		v1.Put("/coupon", updateCoupon)
	}
	//球衣
	{
		v1.Get("/jersey/{id}", getJersey)
		v1.Get("/jerseys", getJerseys)
		v1.Post("/jersey", createJersey)
		v1.Put("/jersey", updateJersey)
	}
	//比赛
	{
		v1.Get("/match/{id}", getMatch)
		v1.Get("/matchs", getMatchs)
		v1.Post("/match", createMatch)
		v1.Put("/match", updateMatch)
	}
	//球队账目
	//球员账目
}
