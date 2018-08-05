package api

import (
	"52jolynn.com/model"
	"52jolynn.com/core"
	"52jolynn.com/mapper"
	"52jolynn.com/misc"
	"fmt"
)

type Mapi interface {
	CreateClub(name string, remark, address, tel *string) *core.Response
}

type mapi struct {
	clubDao mapper.ClubDao
}

func NewMapi(ctx core.Context) Mapi {
	return &mapi{clubDao: mapper.NewClubDao(ctx.Datasource())}
}

//新建俱乐部
func (api *mapi) CreateClub(name string, remark, address, tel *string) *core.Response {
	//判断同名俱乐部是否存在
	club, ok := api.clubDao.GetByName(name)
	if !ok {
		return core.CreateResponse(misc.CodeTryAgainLater)
	}
	if club != nil {
		return core.CreateResponse(misc.CodeDuplicateData, fmt.Sprintf("俱乐部%s", name))
	}

	club = &model.Club{}
	if _, ok := api.clubDao.Insert(club); ok {
		return core.CreateResponse(misc.CodeFailure, "新建俱乐部失败")
	}
	return core.CreateResponseWithData(misc.CodeSuccess, club)
}
