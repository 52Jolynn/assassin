package api

import (
	"52jolynn.com/model"
	"52jolynn.com/core"
	"github.com/kataras/golog"
	"52jolynn.com/mapper"
)

var logger *golog.Logger

type Mapi interface {
	CreateClub(name string, remark, address, tel *string) *model.Club
}

type mapi struct {
	clubDao mapper.ClubDao
}

func NewMapi(ctx core.Context) Mapi {
	logger = ctx.RootLogger().Child("mapi")
	return &mapi{clubDao: mapper.NewClubDao(ctx.Datasource())}
}

//新建俱乐部
func (api *mapi) CreateClub(name string, remark, address, tel *string) *model.Club {
	logger.Info("create club")
	return nil;
}
