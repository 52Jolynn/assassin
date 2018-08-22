package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"52jolynn.com/core"
	"52jolynn.com/misc"
	"strconv"
)

//根据id查询球队
func getTeam(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, uapi.GetGround(id))
}

func getTeams(ctx *gin.Context) {

}

func createTeam(ctx *gin.Context) {
}

func updateTeam(ctx *gin.Context) {

}