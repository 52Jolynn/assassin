package router

import (
	"52jolynn.com/core"
	"52jolynn.com/misc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//根据id查询俱乐部
func getClub(ctx *gin.Context) {
	strId := ctx.Param("id")
	if "" == strId {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "id"))
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamInvalid, "id"))
		return
	}
	ctx.JSON(http.StatusOK, uapi.GetClub(id))
}

func getClubs(ctx *gin.Context) {
	var name, status *string
	if param, exists := ctx.GetQuery("name"); exists {
		name = &param
	}
	if param, exists := ctx.GetQuery("status"); exists {
		status = &param
	}

	var strLimit, strOffset string
	if param, exists := ctx.GetQuery("limit"); exists {
		strLimit = param
	} else {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "limit"))
		return
	}
	if param, exists := ctx.GetQuery("offset"); exists {
		strOffset = param
	} else {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "offset"))
		return
	}
	limit, err := strconv.Atoi(strLimit);
	if err != nil {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamInvalid, "limit"))
		return
	}
	offset, err := strconv.Atoi(strOffset);
	if err != nil {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamInvalid, "offset"))
		return
	}

	ctx.JSON(http.StatusOK, uapi.GetClubs(name, status, limit, offset))
}

func createClub(ctx *gin.Context) {
	name, exists := ctx.GetPostForm("name")
	if !exists {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "name"))
		return
	}

	var remark, address, tel *string
	if param, exists := ctx.GetPostForm("remark"); exists {
		remark = &param
	}
	if param, exists := ctx.GetPostForm("address"); exists {
		address = &param
	}
	if param, exists := ctx.GetPostForm("tel"); exists {
		tel = &param
	}
	ctx.JSON(http.StatusOK, mapi.CreateClub(name, remark, address, tel))
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
