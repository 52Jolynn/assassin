package router

import (
	"52jolynn.com/core"
	"52jolynn.com/api"
	"time"
	"52jolynn.com/misc"
	"github.com/gin-gonic/gin"
	"fmt"
)

var mapi api.Mapi
var uapi api.Uapi

func RegisterRoutes(ctx core.Context, router *gin.Engine) {
	mapi = api.NewMapi(ctx)
	uapi = api.NewUapi(ctx)
	router.GET("/index", func(ctx *gin.Context) {
		ctx.Writer.WriteString(fmt.Sprintf("Hello assassin, today is %s", time.Now().Format(misc.StandardTimeFormatPattern)))
	})
	v1 := router.Group("/v1")
	//俱乐部
	{
		v1.GET("/club/:id", getClub)
		v1.GET("/clubs", getClubs)
		v1.POST("/club", createClub)
		v1.PUT("/club", updateClub)
	}
	//场地
	{
		v1.GET("/ground/:id", getGround)
		v1.GET("/grounds", getGrounds)
		v1.POST("/ground", createGround)
		v1.PUT("/ground", updateGround)
	}
	//球队
	{
		v1.GET("/team/:id", getTeam)
		v1.GET("/teams", getTeams)
		v1.POST("/team", createTeam)
		v1.PUT("/team", updateTeam)
	}
	//球员
	{
		v1.GET("/player/:id", getPlayer)
		v1.GET("/players", getPlayers)
		v1.POST("/player", createPlayer)
		v1.PUT("/player", updatePlayer)
	}
	//优惠券
	{
		v1.GET("/coupon/:id", getCoupon)
		v1.GET("/coupons", getCoupons)
		v1.POST("/coupon", createCoupon)
		v1.PUT("/coupon", updateCoupon)
	}
	//球衣
	{
		v1.GET("/jersey/:id", getJersey)
		v1.GET("/jerseys", getJerseys)
		v1.POST("/jersey", createJersey)
		v1.PUT("/jersey", updateJersey)
	}
	//比赛
	{
		v1.GET("/match/:id", getMatch)
		v1.GET("/matchs", getMatchs)
		v1.POST("/match", createMatch)
		v1.PUT("/match", updateMatch)
	}
	//球队账目
	//球员账目
}
