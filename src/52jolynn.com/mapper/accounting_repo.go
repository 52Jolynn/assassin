package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type AccountingOfTeamDao interface {
	InsertAccountingOfTeam(accounting *model.AccountingOfTeam) (*model.AccountingOfTeam, bool)
	QueryAccountingOfTeam(matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString, limit, offset int) ([]model.AccountingOfTeam, bool)
}

type accountingOfTeamDao struct {
	db *sql.DB
}

func NewAccountingOfTeamDao(db *sql.DB) AccountingOfTeamDao {
	return &accountingOfTeamDao{db: db}
}

const (
	ColumnWithoutIdOfAccountingOfTeam = "`ref_id`, `match_id`, `team_id`, `amount`, `remark`, `before_balance`, `after_balance`, `type`, `bill_date`, `create_time`"
	ColumnAccountingOfTeam            = "`id`, " + ColumnWithoutIdOfAccountingOfTeam
	TableNameOfAccountingOfTeam       = "accounting_of_team"
)

func (a *accountingOfTeamDao) InsertAccountingOfTeam(accounting *model.AccountingOfTeam) (*model.AccountingOfTeam, bool) {
	stmt, err := a.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", TableNameOfAccountingOfTeam, ColumnAccountingOfTeam))
	if err != nil {
		log.Printf("预编译插入%s语句出错，err: %s\n", TableNameOfAccountingOfTeam, err.Error())
		return nil, false
	}
	defer stmt.Close()
	_, err = stmt.Exec(accounting.Id, accounting.RefId, accounting.MatchId, accounting.TeamId, accounting.Amount, accounting.Remark, accounting.BeforeBalance, accounting.AfterBalance, accounting.Ttype, accounting.BillDate, accounting.CreateTime)
	if err != nil {
		log.Printf("插入%s出错，err: %s\n", TableNameOfAccountingOfTeam, err.Error())
		return nil, false
	}
	return accounting, true
}

func (a *accountingOfTeamDao) QueryAccountingOfTeam(matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString, limit, offset int) ([]model.AccountingOfTeam, bool) {
	querySql, args := buildQueryAccountingOfTeam(ColumnPlayerStatOfMatch, matchId, teamId, ttype, startTime, endTime)
	querySql += " order by create_time desc, id"
	return a.queryAccoutingOfTeam(querySql, args...)
}

func (a *accountingOfTeamDao) queryAccoutingOfTeam(query string, args ...interface{}) ([]model.AccountingOfTeam, bool) {
	stmt, err := a.db.Prepare(query)
	sqlMsg := fmt.Sprintf("%s.queryAccoutingOfTeam", TableNameOfAccountingOfTeam)
	if err != nil {
		log.Printf("预编译%s语句出错，err: %s\n", sqlMsg, err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("%s查询出错，err: %s\n", sqlMsg, err.Error())
		return nil, false
	}

	accountingOfTeams := make([]model.AccountingOfTeam, 0)
	for rows.Next() {
		accounting := model.AccountingOfTeam{}
		err = rows.Scan(&accounting.Id, &accounting.RefId, &accounting.MatchId, &accounting.TeamId, &accounting.Amount, &accounting.Remark, &accounting.BeforeBalance, &accounting.AfterBalance, &accounting.Ttype, &accounting.BillDate, &accounting.CreateTime)
		if err != nil {
			log.Printf("%s获取数据出错，err: %s\n", sqlMsg, err.Error())
			return nil, false
		}
		accountingOfTeams = append(accountingOfTeams, accounting)
	}

	return accountingOfTeams, true
}

func buildQueryAccountingOfTeam(returnColumn string, matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfAccountingOfTeam))
	var args []interface{}
	querySql.WriteString(" and match_id=?")
	args = append(args, matchId)
	if matchId.Valid {
		querySql.WriteString(" and match_id=?")
		args = append(args, matchId)
	}
	if teamId.Valid {
		querySql.WriteString(" and team_id=?")
		args = append(args, teamId)
	}
	if ttype.Valid {
		querySql.WriteString(" and type=?")
		args = append(args, ttype)
	}
	if startTime.Valid {
		querySql.WriteString(" and create_time>=?")
		args = append(args, startTime)
	}
	if endTime.Valid {
		querySql.WriteString(" and create_time<=?")
		args = append(args, endTime)
	}
	return querySql.String(), args
}

type AccountingOfPlayerDao interface {
	InsertAccountingOfPlayer(accounting *model.AccountingOfPlayer) (*model.AccountingOfPlayer, bool)
	QueryAccountingOfPlayer(matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString, limit, offset int) ([]model.AccountingOfPlayer, bool)
}

type accountingOfPlayerDao struct {
	db *sql.DB
}

func NewAccountingOfPlayerDao(db *sql.DB) AccountingOfPlayerDao {
	return &accountingOfPlayerDao{db: db}
}

const (
	ColumnWithoutIdOfAccountingOfPlayer = "`ref_id`, `match_id`, `team_id`, `player_id`, `amount`, `remark`, `before_balance`, `after_balance`, `bill_date`, `type`, `create_time`"
	ColumnAccountingOfPlayer            = "`id`, " + ColumnWithoutIdOfAccountingOfPlayer
	TableNameOfAccountingOfPlayer       = "accounting_of_player"
)

func (a *accountingOfPlayerDao) InsertAccountingOfPlayer(accounting *model.AccountingOfPlayer) (*model.AccountingOfPlayer, bool) {
	stmt, err := a.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", TableNameOfAccountingOfPlayer, ColumnAccountingOfPlayer))
	if err != nil {
		log.Printf("预编译插入%s语句出错，err: %s\n", TableNameOfAccountingOfPlayer, err.Error())
		return nil, false
	}
	defer stmt.Close()
	_, err = stmt.Exec(accounting.Id, accounting.RefId, accounting.MatchId, accounting.TeamId, accounting.PlayerId, accounting.Amount, accounting.Remark, accounting.BeforeBalance, accounting.AfterBalance, accounting.Ttype, accounting.BillDate, accounting.CreateTime)
	if err != nil {
		log.Printf("插入%s出错，err: %s\n", TableNameOfAccountingOfPlayer, err.Error())
		return nil, false
	}
	return accounting, true
}

func (a *accountingOfPlayerDao) QueryAccountingOfPlayer(matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString, limit, offset int) ([]model.AccountingOfPlayer, bool) {
	querySql, args := buildQueryAccountingOfPlayer(ColumnPlayerStatOfMatch, matchId, teamId, ttype, startTime, endTime)
	querySql += " order by create_time desc, id"
	return a.queryAccoutingOfPlayer(querySql, args...)
}

func (a *accountingOfPlayerDao) queryAccoutingOfPlayer(query string, args ...interface{}) ([]model.AccountingOfPlayer, bool) {
	stmt, err := a.db.Prepare(query)
	sqlMsg := fmt.Sprintf("%s.queryAccoutingOfPlayer", TableNameOfAccountingOfPlayer)
	if err != nil {
		log.Printf("预编译%s语句出错，err: %s\n", sqlMsg, err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("%s查询出错，err: %s\n", sqlMsg, err.Error())
		return nil, false
	}

	accountingOfPlayers := make([]model.AccountingOfPlayer, 0)
	for rows.Next() {
		accounting := model.AccountingOfPlayer{}
		err = rows.Scan(&accounting.Id, &accounting.RefId, &accounting.MatchId, &accounting.TeamId, &accounting.PlayerId, &accounting.Amount, &accounting.Remark, &accounting.BeforeBalance, &accounting.AfterBalance, &accounting.Ttype, &accounting.BillDate, &accounting.CreateTime)
		if err != nil {
			log.Printf("%s获取数据出错，err: %s\n", sqlMsg, err.Error())
			return nil, false
		}
		accountingOfPlayers = append(accountingOfPlayers, accounting)
	}

	return accountingOfPlayers, true
}

func buildQueryAccountingOfPlayer(returnColumn string, matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfAccountingOfPlayer))
	var args []interface{}
	querySql.WriteString(" and match_id=?")
	args = append(args, matchId)
	if matchId.Valid {
		querySql.WriteString(" and match_id=?")
		args = append(args, matchId)
	}
	if teamId.Valid {
		querySql.WriteString(" and team_id=?")
		args = append(args, teamId)
	}
	if ttype.Valid {
		querySql.WriteString(" and type=?")
		args = append(args, ttype)
	}
	if startTime.Valid {
		querySql.WriteString(" and create_time>=?")
		args = append(args, startTime)
	}
	if endTime.Valid {
		querySql.WriteString(" and create_time<=?")
		args = append(args, endTime)
	}
	return querySql.String(), args
}
