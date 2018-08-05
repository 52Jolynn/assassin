package api

import (
	"52jolynn.com/mapper"
	"52jolynn.com/core"
	"52jolynn.com/misc"
	"fmt"
)

type Uapi interface {
	getClub(id int) *core.Response
	getClubs(limit, offset int) *core.Response
}

type uapi struct {
	clubDao mapper.ClubDao
}

func NewUapi(ctx core.Context) Uapi {
	return &uapi{clubDao: mapper.NewClubDao(ctx.Datasource())}
}

func (u *uapi) getClub(id int) *core.Response {
	club, ok := u.clubDao.GetById(id)
	if !ok {
		return core.CreateResponse(misc.CodeDataDoesNotExist, fmt.Sprintf("俱乐部%s", id))
	}
	return core.CreateResponse(misc.CodeSuccess, club)

}

func (u *uapi) getClubs(limit, offset int) *core.Response {
	return nil
}
