package api

import (
	"52jolynn.com/model"
	"52jolynn.com/core"
	"52jolynn.com/mapper"
	"52jolynn.com/misc"
	"fmt"
	"time"
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

const (
	ClubStatusNormal  = "N"
	ClubStatusDisable = "D"
)

//新建俱乐部
func (api *mapi) CreateClub(name string, remark, address, tel *string) *core.Response {
	//判断同名俱乐部是否存在
	ok, exists := api.clubDao.ExistsByName(name)
	if !ok {
		return core.CreateResponse(misc.CodeTryAgainLater)
	}
	if exists {
		return core.CreateResponse(misc.CodeDuplicateData, fmt.Sprintf("俱乐部%s", name))
	}

	club := &model.Club{Name: name, Remark: remark, Address: address, Tel: tel, CreateTime: time.Now().Format(misc.StandardTimeFormatPattern), Status: ClubStatusNormal}
	if _, ok := api.clubDao.Insert(club); !ok {
		return core.CreateResponse(misc.CodeFailure, "新建俱乐部失败")
	}
	return core.CreateResponseWithData(misc.CodeSuccess, club)
}
