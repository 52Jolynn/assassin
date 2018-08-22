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
	if param, exists := ctx.GetPostForm("name"); exists {
		name = &param
	}
	if param, exists := ctx.GetPostForm("status"); exists {
		status = &param
	}

	var strLimit, strOffset string
	if param, exists := ctx.GetPostForm("limit"); exists {
		strLimit = param
	} else {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "limit"))
		return
	}
	if param, exists := ctx.GetPostForm("offset"); exists {
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

	ctx.JSON(http.StatusOK, uapi.QueryClub(name, status, limit, offset))
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

	var name, remark, address, tel *string
	if param, exists := ctx.GetPostForm("name"); exists {
		name = &param
	}
	if param, exists := ctx.GetPostForm("remark"); exists {
		remark = &param
	}
	if param, exists := ctx.GetPostForm("address"); exists {
		address = &param
	}
	if param, exists := ctx.GetPostForm("tel"); exists {
		tel = &param
	}
	ctx.JSON(http.StatusOK, mapi.UpdateClub(id, name, remark, address, tel))
}

//根据id查询场地
func getGround(ctx *gin.Context) {
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

func getGrounds(ctx *gin.Context) {
	var name, ttype, status *string
	var clubId *int
	if param, exists := ctx.GetPostForm("name"); exists {
		name = &param
	}
	if param, exists := ctx.GetPostForm("type"); exists {
		ttype = &param
	}
	if param, exists := ctx.GetPostForm("status"); exists {
		status = &param
	}
	if strClubId, exists := ctx.GetPostForm("club_id"); exists {
		value, err := strconv.Atoi(strClubId)
		if err != nil {
			ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamInvalid, "club_id"))
			return
		}
		clubId = &value
	}

	var strLimit, strOffset string
	if param, exists := ctx.GetPostForm("limit"); exists {
		strLimit = param
	} else {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "limit"))
		return
	}
	if param, exists := ctx.GetPostForm("offset"); exists {
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

	ctx.JSON(http.StatusOK, uapi.QueryGround(name, ttype, status, clubId, limit, offset))
}

func createGround(ctx *gin.Context) {
	name, exists := ctx.GetPostForm("name")
	if !exists {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "name"))
		return
	}

	ttype, exists := ctx.GetPostForm("type")
	if !exists {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "type"))
		return
	}

	strClubId := ctx.Param("club_id")
	if "" == strClubId {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamMissing, "club_id"))
		return
	}

	clubId, err := strconv.Atoi(strClubId)
	if err != nil {
		ctx.JSON(http.StatusOK, core.CreateResponse(misc.CodeParamInvalid, "club_id"))
		return
	}

	var remark *string
	if param, exists := ctx.GetPostForm("remark"); exists {
		remark = &param
	}
	ctx.JSON(http.StatusOK, mapi.CreateGround(name, ttype, clubId, remark))
}

func updateGround(ctx *gin.Context) {
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
	var name, remark *string
	if param, exists := ctx.GetPostForm("name"); exists {
		name = &param
	}
	if param, exists := ctx.GetPostForm("remark"); exists {
		remark = &param
	}

	ctx.JSON(http.StatusOK, mapi.UpdateGround(id, name, remark))

}
