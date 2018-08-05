package api

import (
	"52jolynn.com/model"
	"52jolynn.com/core"
	"github.com/kataras/golog"
	"52jolynn.com/mapper"
)

var mlogger *golog.Logger

type Mapi interface {
	CreateClub(name string, remark, address, tel *string) *model.Club
}

type mapi struct {
	clubDao mapper.ClubDao
}

func NewMapi(ctx core.Context) Mapi {
	mlogger = ctx.RootLogger().Child("mapi")
	return &mapi{clubDao: mapper.NewClubDao(ctx.Datasource())}
}

//新建俱乐部
func (api *mapi) CreateClub(name string, remark, address, tel *string) *model.Club {
	mlogger.Info("create club")
	return &model.Club{};
}
