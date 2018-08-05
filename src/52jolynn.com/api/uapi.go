package api

import (
	"github.com/kataras/golog"
	"52jolynn.com/model"
	"52jolynn.com/mapper"
	"52jolynn.com/core"
)

var ulogger *golog.Logger

type Uapi interface {
	getClub(id int) *model.Club
	getClubs(limit, offset int) *core.Pagination
}

type uapi struct {
	clubDao mapper.ClubDao
}

func NewUapi(ctx core.Context) Uapi {
	ulogger = ctx.RootLogger().Child("uapi")
	return &uapi{clubDao: mapper.NewClubDao(ctx.Datasource())}
}

func (u *uapi) getClub(id int) *model.Club {
	ulogger.Info("get club %d", id)
	return nil
}

func (u *uapi) getClubs(limit, offset int) *core.Pagination {
	return nil
}
