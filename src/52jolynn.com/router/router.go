package router

import (
	"52jolynn.com/core"
	"52jolynn.com/api"
	"time"
	"52jolynn.com/misc"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

var mapi api.Mapi
var uapi api.Uapi

var session = sessions.NewSession(sessions.NewCookieStore([]byte("assassin-secret")), "gin-session")

const (
	SessionUid = "uid"
)

func RegisterRoutes(ctx core.Context, router *gin.Engine) {
	mapi = api.NewMapi(ctx)
	uapi = api.NewUapi(ctx)

	router.GET("/index", func(ctx *gin.Context) {
		ctx.Writer.WriteString(fmt.Sprintf("Hello assassin, today is %s", time.Now().Format(misc.StandardTimeFormatPattern)))
	})

	v1 := router.Group("/v1")
	//用户
	{
		v1.POST("/user/sign_up", signUp)
		v1.POST("/user/sign_in", signIn)
		v1.POST("/user/sign_with_wx", signWithWeiXin)
	}

	authV1 := router.Group("/v1")
	authV1.Use(func(ctx *gin.Context) {
		if _, exists := session.Values[SessionUid]; !exists {
			ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeUserDoesNotSignIn))
			ctx.Abort()
			return
		}
		ctx.Next()
	})
	//用户
	{
		authV1.POST("/user/sign_out", signOut)
		authV1.POST("/user/bind_wx", bindWeiXin)
	}
	//俱乐部
	{
		authV1.GET("/club/:id", getClub)
		authV1.GET("/clubs", getClubs)
		authV1.POST("/club", createClub)
		authV1.PUT("/club", updateClub)
	}
	//场地
	{
		authV1.GET("/ground/:id", getGround)
		authV1.GET("/grounds", getGrounds)
		authV1.POST("/ground", createGround)
		authV1.PUT("/ground", updateGround)
	}
	//球队
	{
		authV1.GET("/team/:id", getTeam)
		authV1.GET("/teams", getTeams)
		authV1.POST("/team", createTeam)
		authV1.PUT("/team", updateTeam)
	}
	//球员
	{
		authV1.GET("/player/:id", getPlayer)
		authV1.GET("/players", getPlayers)
		authV1.POST("/player", createPlayer)
		authV1.PUT("/player", updatePlayer)
	}
	//优惠券
	{
		authV1.GET("/coupon/:id", getCoupon)
		authV1.GET("/coupons", getCoupons)
		authV1.POST("/coupon", createCoupon)
	}
	//球衣
	{
		authV1.GET("/jersey/:id", getJersey)
		authV1.GET("/jerseys", getJerseys)
		authV1.POST("/jersey", createJersey)
		authV1.PUT("/jersey", updateJersey)
	}
	//比赛
	{
		authV1.GET("/match/:id", getMatch)
		authV1.GET("/matchs", getMatchs)
		authV1.POST("/match", createMatch)
		authV1.PUT("/match", updateMatch)
	}
	//球队账目
	//球员账目
}
