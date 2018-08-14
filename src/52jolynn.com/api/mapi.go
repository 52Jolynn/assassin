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
	UpdateClub(id int, name, remark, address, tel *string) *core.Response

	CreateGround(name string, ttype string, clubId int, remark *string) *core.Response
	UpdateGround(id int, name, remark *string) *core.Response
}

type mapi struct {
	clubDao   mapper.ClubDao
	groundDao mapper.GroundDao
}

func NewMapi(ctx core.Context) Mapi {
	return &mapi{clubDao: mapper.NewClubDao(ctx.Datasource()), groundDao: mapper.NewGroundDao(ctx.Datasource())}
}

const (
	ClubStatusNormal  = "N"
	ClubStatusDisable = "D"

	GroundStatusNormal  = "N"
	GroundStatusDisable = "D"
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

//更新俱乐部信息
func (api *mapi) UpdateClub(id int, name, remark, address, tel *string) *core.Response {
	club, ok := api.clubDao.GetById(id)
	if !ok {
		return core.CreateResponse(misc.CodeDataDoesNotExist, fmt.Sprintf("%d", id))
	}
	if name != nil {
		club.Name = *name
	}
	if address != nil {
		club.Address = address
	}
	if remark != nil {
		club.Remark = remark
	}
	if tel != nil {
		club.Tel = tel
	}
	if effectRows, ok := api.clubDao.Update(club); !ok || effectRows != 1 {
		return core.CreateResponse(misc.CodeFailure, "更新俱乐部信息失败")
	}
	return core.CreateResponse(misc.CodeSuccess)
}

func (api *mapi) CreateGround(name string, ttype string, clubId int, remark *string) *core.Response {
	ground := &model.Ground{Name: name, Remark: remark, Ttype: ttype, ClubId: clubId, CreateTime: time.Now().Format(misc.StandardTimeFormatPattern), Status: ClubStatusNormal}
	if _, ok := api.groundDao.Insert(ground); !ok {
		return core.CreateResponse(misc.CodeFailure, "新建俱乐部场地失败")
	}
	return core.CreateResponseWithData(misc.CodeSuccess, ground)
}

func (api *mapi) UpdateGround(id int, name, remark *string) *core.Response {
	ground, ok := api.groundDao.GetById(id)
	if !ok {
		return core.CreateResponse(misc.CodeDataDoesNotExist, fmt.Sprintf("%d", id))
	}
	if name != nil {
		ground.Name = *name
	}
	if remark != nil {
		ground.Remark = remark
	}
	if effectRows, ok := api.groundDao.Update(ground); !ok || effectRows != 1 {
		return core.CreateResponse(misc.CodeFailure, "更新俱乐部场地信息失败")
	}
	return core.CreateResponse(misc.CodeSuccess)
}
