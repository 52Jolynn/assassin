package route

import (
	"github.com/kataras/iris"
	"52jolynn.com/core"
	"52jolynn.com/misc"
)

//根据id查询俱乐部
func getClub(ctx iris.Context) {

}

func getClubs(ctx iris.Context) {

}

func createClub(ctx iris.Context) {
	club := mapi.CreateClub(ctx.PostValue("name"), nil, nil, nil)
	ctx.JSON(core.ResponseWithData{Response: core.Response{Code: misc.CodeSuccess, Msg: misc.ResponseCode[misc.CodeSuccess]}, Data: club})
}

func updateClub(ctx iris.Context) {

}

//根据id查询场地
func getGround(ctx iris.Context) {

}

func getGrounds(ctx iris.Context) {

}

func createGround(ctx iris.Context) {
	club := mapi.CreateClub(ctx.PostValue("name"), nil, nil, nil)
	ctx.JSON(core.ResponseWithData{Response: core.Response{Code: misc.CodeSuccess, Msg: misc.ResponseCode[misc.CodeSuccess]}, Data: club})
}

func updateGround(ctx iris.Context) {

}

//根据id查询球队
func getTeam(ctx iris.Context) {

}

func getTeams(ctx iris.Context) {

}

func createTeam(ctx iris.Context) {
	club := mapi.CreateClub(ctx.PostValue("name"), nil, nil, nil)
	ctx.JSON(core.ResponseWithData{Response: core.Response{Code: misc.CodeSuccess, Msg: misc.ResponseCode[misc.CodeSuccess]}, Data: club})
}

func updateTeam(ctx iris.Context) {

}

//根据id查询球员
func getPlayer(ctx iris.Context) {

}

func getPlayers(ctx iris.Context) {

}

func createPlayer(ctx iris.Context) {
	club := mapi.CreateClub(ctx.PostValue("name"), nil, nil, nil)
	ctx.JSON(core.ResponseWithData{Response: core.Response{Code: misc.CodeSuccess, Msg: misc.ResponseCode[misc.CodeSuccess]}, Data: club})
}

func updatePlayer(ctx iris.Context) {

}

//根据id查询优惠券
func getCoupon(ctx iris.Context) {

}

func getCoupons(ctx iris.Context) {

}

func createCoupon(ctx iris.Context) {
	club := mapi.CreateClub(ctx.PostValue("name"), nil, nil, nil)
	ctx.JSON(core.ResponseWithData{Response: core.Response{Code: misc.CodeSuccess, Msg: misc.ResponseCode[misc.CodeSuccess]}, Data: club})
}

func updateCoupon(ctx iris.Context) {

}

//根据id查询球衣
func getJersey(ctx iris.Context) {

}

func getJerseys(ctx iris.Context) {

}

func createJersey(ctx iris.Context) {
	club := mapi.CreateClub(ctx.PostValue("name"), nil, nil, nil)
	ctx.JSON(core.ResponseWithData{Response: core.Response{Code: misc.CodeSuccess, Msg: misc.ResponseCode[misc.CodeSuccess]}, Data: club})
}

func updateJersey(ctx iris.Context) {

}

//根据id查询比赛
func getMatch(ctx iris.Context) {

}

func getMatchs(ctx iris.Context) {

}

func createMatch(ctx iris.Context) {
	club := mapi.CreateClub(ctx.PostValue("name"), nil, nil, nil)
	ctx.JSON(core.ResponseWithData{Response: core.Response{Code: misc.CodeSuccess, Msg: misc.ResponseCode[misc.CodeSuccess]}, Data: club})
}

func updateMatch(ctx iris.Context) {

}
