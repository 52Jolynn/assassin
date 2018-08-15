package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
	"fmt"
	"strings"
)

type TeamDao interface {
	GetById(id int) (*model.Team, bool)
	Insert(team *model.Team) (*model.Team, bool)
	QueryTeam(name *string, status []string, limit, offset int) (*[]model.Team, bool)
	QueryCount(name *string, status []string) int
	Update(team *model.Team) (int64, bool)

	JoinClub(tc *model.TeamOfClub) (*model.TeamOfClub, bool)

	GetJerseyById(id int) (*model.JerseyOfTeam, bool)
	QueryJersey(teamId int, status []string) (*[]model.JerseyOfTeam, bool)
	InsertJersey(j *model.JerseyOfTeam) (*model.JerseyOfTeam, bool)
	UpdateJerseyStatus(id int, status, oldStatus string) (int64, bool)
}

type teamDao struct {
	db *sql.DB
}

func NewTeamDao(db *sql.DB) TeamDao {
	return &teamDao{db: db}
}

const (
	ColumnWithoutIdOfTeam = "`name`, `remark`, `captain_name`, `captain_mobile`, `manager_username`, `manager_passwd`, `create_time`, `status`"
	ColumnOfTeam          = "`id`, " + ColumnWithoutIdOfTeam
	TableNameOfTeam       = "team"
)

//根据id获取球队信息
func (c *teamDao) GetById(id int) (*model.Team, bool) {
	if teams, ok := c.queryTeam(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfTeam, TableNameOfTeam), id); ok && len(*teams) == 1 {
		return &(*teams)[0], true
	}
	return nil, false
}

func buildQueryteamSql(returnColumn string, name *string, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfTeam))
	var args []interface{}
	if name != nil {
		querySql.WriteString(" and name=?")
		args = append(args, name)
	}
	statusLen := len(status)
	if statusLen > 0 {
		querySql.WriteString(" and `status` in(")
		for index, value := range status {
			querySql.WriteString("?")
			if index != statusLen-1 {
				querySql.WriteString(",")
			}
			args = append(args, value)
		}
		querySql.WriteString(")")
	}
	return querySql.String(), args
}

//搜索
func (c *teamDao) QueryTeam(name *string, status []string, limit, offset int) (*[]model.Team, bool) {
	querySql, args := buildQueryteamSql(ColumnOfTeam, name, status)
	querySql += " order by create_time desc, id limit ? offset ?"
	args = append(args, limit, offset)
	return c.queryTeam(querySql, args...)
}

//搜索计数
func (c *teamDao) QueryCount(name *string, status []string) int {
	querySql, args := buildQueryteamSql("count(*)", name, status)
	stmt, err := c.db.Prepare(querySql)
	if err != nil {
		log.Printf("预编译teamDao.QueryCount语句出错，err: %s", err.Error())
		return 0
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("teamDao.QueryCount查询出错，err: %s", err.Error())
		return 0
	}
	if !rows.Next() {
		return 0
	}
	count := 0
	err = rows.Scan(&count)
	if err != nil {
		log.Printf("teamDao.QueryCount获取数据出错，err: %s", err.Error())
	}
	return count
}

func (c *teamDao) queryTeam(query string, args ...interface{}) (*[]model.Team, bool) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Printf("预编译teamDao.queryTeam语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("teamDao.queryTeam查询出错，err: %s", err.Error())
		return nil, false
	}

	teams := make([]model.Team, 0)
	for rows.Next() {
		team := model.Team{}
		err = rows.Scan(&team.Id, &team.Name, &team.Remark, &team.CaptainName, &team.CaptainMobile, &team.ManagerUsername, &team.ManagerPasswd, &team.CreateTime, &team.Status)
		if err != nil {
			log.Printf("teamDao.queryTeam获取数据出错，err: %s", err.Error())
			return nil, false
		}
		teams = append(teams, team)
	}

	return &teams, true
}

func (c *teamDao) Insert(team *model.Team) (*model.Team, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?, ?, ?)", TableNameOfTeam, ColumnWithoutIdOfTeam))
	if err != nil {
		log.Printf("预编译插入team语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(team.Name, team.Remark, team.CaptainName, team.CaptainMobile, team.ManagerUsername, team.ManagerPasswd, team.CreateTime, team.Status)
	if err != nil {
		log.Printf("插入team出错，err: %s", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入team.id出错，err: %s", err.Error())
		return nil, false
	}
	team.Id = int(lastInsertId)
	return team, true
}

func (c *teamDao) Update(team *model.Team) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `name`=?, `remark`=?, `captain_name`=?, `captain_mobile`=?, `manager_username`=?, `manager_passwd`=?, `create_time`=?, `status`=? where id=?", TableNameOfTeam))
	if err != nil {
		log.Printf("预编译更新team语句出错，err: %s", err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(team.Name, team.Remark, team.CaptainName, team.CaptainMobile, team.ManagerUsername, team.ManagerPasswd, team.CreateTime, team.Status, team.Id)
	if err != nil {
		log.Printf("更新team出错，err: %s", err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新team影响行数出错，err: %s", err.Error())
		return 0, false
	}
	return rowsAffected, true
}

const (
	ColumnWithoutIdOfTeamOfClub = "`club_id`, `team_id`, `present_balance`, `used_amount`, `join_time`, `create_time`"
	TableNameOfTeamOfClub       = "team_of_club"
)

func (c *teamDao) JoinClub(tc *model.TeamOfClub) (*model.TeamOfClub, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?)", TableNameOfTeamOfClub, ColumnWithoutIdOfTeamOfClub))
	if err != nil {
		log.Printf("预编译插入team_of_club语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(tc.ClubId, tc.TeamId, tc.PresentBalance, tc.UsedAmount, tc.JoinTime, tc.CreateTime)
	if err != nil {
		log.Printf("插入team_of_club出错，err: %s", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入team_of_club.id出错，err: %s", err.Error())
		return nil, false
	}
	tc.Id = int(lastInsertId)
	return tc, true
}

const (
	ColumnWithoutIdOfJerseyOfTeam = "`team_id`, `home_color`, `away_color`, `create_time`, `status`"
	ColumnOfJerseyOfTeam          = "`id`, " + ColumnWithoutIdOfJerseyOfTeam
	TableNameOfJerseyOfTeam       = "jersey_of_team"
)

func (c *teamDao) GetJerseyById(id int) (*model.JerseyOfTeam, bool) {
	if jerseys, ok := c.queryJerseyOfTeam(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfJerseyOfTeam, TableNameOfJerseyOfTeam), id); ok && len(*jerseys) == 1 {
		return &(*jerseys)[0], true
	}
	return nil, false
}

func (c *teamDao) QueryJersey(teamId int, status []string) (*[]model.JerseyOfTeam, bool) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where team_id=?", ColumnOfJerseyOfTeam, TableNameOfJerseyOfTeam))

	var args []interface{}
	args = append(args, teamId)

	statusLen := len(status)
	if statusLen > 0 {
		querySql.WriteString(" and `status` in(")
		for index, value := range status {
			querySql.WriteString("?")
			if index != statusLen-1 {
				querySql.WriteString(",")
			}
			args = append(args, value)
		}
		querySql.WriteString(")")
	}
	return c.queryJerseyOfTeam(querySql.String(), args...)
}

func (c *teamDao) InsertJersey(jersey *model.JerseyOfTeam) (*model.JerseyOfTeam, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?)", TableNameOfJerseyOfTeam, ColumnWithoutIdOfJerseyOfTeam))
	if err != nil {
		log.Printf("预编译插入jersey_of_team语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(jersey.TeamId, jersey.HomeColor, jersey.AwayColor, jersey.CreateTime, jersey.Status)
	if err != nil {
		log.Printf("插入jersey_of_team出错，err: %s", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入jersey_of_team.id出错，err: %s", err.Error())
		return nil, false
	}
	jersey.Id = int(lastInsertId)
	return jersey, true
}

func (c *teamDao) UpdateJerseyStatus(id int, status, oldStatus string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update %s set `status`=? where id=? and `status`=?", TableNameOfJerseyOfTeam))
	if err != nil {
		log.Printf("预编译更新jersey_of_team.status语句出错，err: %s", err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(status, id, oldStatus)
	if err != nil {
		log.Printf("更新jersey_of_team.status出错，err: %s", err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新jersey_of_team.status影响行数出错，err: %s", err.Error())
		return 0, false
	}
	return rowsAffected, true
}

func (c *teamDao) queryJerseyOfTeam(query string, args ...interface{}) (*[]model.JerseyOfTeam, bool) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Printf("预编译teamDao.queryJerseyOfTeam语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("teamDao.queryJerseyOfTeam查询出错，err: %s", err.Error())
		return nil, false
	}

	jerseys := make([]model.JerseyOfTeam, 0)
	for rows.Next() {
		jersey := model.JerseyOfTeam{}
		err = rows.Scan(&jersey.Id, &jersey.TeamId, &jersey.HomeColor, &jersey.AwayColor, &jersey.CreateTime, &jersey.Status)
		if err != nil {
			log.Printf("teamDao.queryJerseyOfTeam获取数据出错，err: %s", err.Error())
			return nil, false
		}
		jerseys = append(jerseys, jersey)
	}

	return &jerseys, true
}
