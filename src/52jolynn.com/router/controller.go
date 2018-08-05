package router

import (
	"52jolynn.com/core"
	"52jolynn.com/misc"
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
)

//根据id查询俱乐部
func getClub(ctx *gin.Context) {

}

func getClubs(ctx *gin.Context) {

}

func createClub(ctx *gin.Context) {
	name, exists := ctx.GetPostForm("name")
	if !exists {
		ctx.JSON(http.StatusOK, core.Response{Code: misc.CodeParamMissing, Msg: fmt.Sprintf(misc.ResponseCode[misc.CodeParamMissing], "name")})
	}
	remark, exists := ctx.GetPostForm("remark")
	address, exists := ctx.GetPostForm("address")
	tel, exists := ctx.GetPostForm("tel")

	ctx.JSON(http.StatusOK, mapi.CreateClub(name, &remark, &address, &tel))
}

func updateClub(ctx *gin.Context) {

}

//根据id查询场地
func getGround(ctx *gin.Context) {

}

func getGrounds(ctx *gin.Context) {

}

func createGround(ctx *gin.Context) {
}

func updateGround(ctx *gin.Context) {

}

//根据id查询球队
func getTeam(ctx *gin.Context) {

}

func getTeams(ctx *gin.Context) {

}

func createTeam(ctx *gin.Context) {
}

func updateTeam(ctx *gin.Context) {

}

//根据id查询球员
func getPlayer(ctx *gin.Context) {

}

func getPlayers(ctx *gin.Context) {

}

func createPlayer(ctx *gin.Context) {
}

func updatePlayer(ctx *gin.Context) {

}

//根据id查询优惠券
func getCoupon(ctx *gin.Context) {

}

func getCoupons(ctx *gin.Context) {

}

func createCoupon(ctx *gin.Context) {
}

func updateCoupon(ctx *gin.Context) {

}

//根据id查询球衣
func getJersey(ctx *gin.Context) {

}

func getJerseys(ctx *gin.Context) {

}

func createJersey(ctx *gin.Context) {
}

func updateJersey(ctx *gin.Context) {

}

//根据id查询比赛
func getMatch(ctx *gin.Context) {

}

func getMatchs(ctx *gin.Context) {

}

func createMatch(ctx *gin.Context) {
}

func updateMatch(ctx *gin.Context) {

}
