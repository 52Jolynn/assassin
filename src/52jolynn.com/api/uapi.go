package api

import (
	"52jolynn.com/mapper"
	"52jolynn.com/core"
	"52jolynn.com/misc"
	"fmt"
)

type Uapi interface {
	GetClub(id int) *core.Response
	GetClubs(name *string, status *string, limit, offset int) *core.Response
}

type uapi struct {
	clubDao mapper.ClubDao
}

func NewUapi(ctx core.Context) Uapi {
	return &uapi{clubDao: mapper.NewClubDao(ctx.Datasource())}
}

func (u *uapi) GetClub(id int) *core.Response {
	club, ok := u.clubDao.GetById(id)
	if !ok {
		return core.CreateResponse(misc.CodeDataDoesNotExist, fmt.Sprintf("俱乐部%s", id))
	}
	return core.CreateResponse(misc.CodeSuccess, club)

}

func (u *uapi) GetClubs(name *string, status *string, limit, offset int) *core.Response {
	var statuses []string
	if status == nil {
		statuses = append(statuses, ClubStatusNormal, ClubStatusDisable)
	} else {
		statuses = append(statuses, *status)
	}
	result, ok := u.clubDao.QueryClub(name, statuses, limit, offset)
	if !ok {
		return core.CreateResponse(misc.CodeTryAgainLater)
	}

	return core.CreateResponseWithData(misc.CodeSuccess, core.NewPagination(limit, offset, u.clubDao.QueryCount(name, statuses), result))
}
