package api

import (
		"52jolynn.com/model"
	"52jolynn.com/mapper"
	"52jolynn.com/core"
)

type Uapi interface {
	getClub(id int) *model.Club
	getClubs(limit, offset int) *core.Pagination
}

type uapi struct {
	clubDao mapper.ClubDao
}

func NewUapi(ctx core.Context) Uapi {
	return &uapi{clubDao: mapper.NewClubDao(ctx.Datasource())}
}

func (u *uapi) getClub(id int) *model.Club {
	return nil
}

func (u *uapi) getClubs(limit, offset int) *core.Pagination {
	return nil
}
