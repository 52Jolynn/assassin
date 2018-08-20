package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
	"fmt"
	"strings"
)

type GameOfMatchDao interface {
	GetById(id int) (*model.GameOfMatch, bool)
	Insert(gameOfMatch *model.GameOfMatch) (*model.GameOfMatch, bool)
	QueryGameOfMatch(clubId sql.NullInt64, startOpenTime, endOpenTime sql.NullString, status []string, limit, offset int) ([]model.GameOfMatch, bool)
	QueryCount(clubId sql.NullInt64, startOpenTime, endOpenTime sql.NullString, status []string) int
	UpdateStatus(id int, status, oldStatus string) (int64, bool)
	UpdateEnrollCount(id, count int) (int64, bool)

	QueryEnrollOfMatch(matchId int, playerId sql.NullInt64, status []string) ([]model.EnrollOfMatch, bool)
	InsertEnrollOfMatch(e *model.EnrollOfMatch) (*model.EnrollOfMatch, bool)
	UpdateEnrollOfMatchStatus(id int64, status, oldStatus string) (int64, bool)

	InsertTeamStatOfMatch(stat *model.TeamStatOfMatch) (*model.TeamStatOfMatch, bool)
	QueryTeamStatOfMatch(matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString, limit, offset int) ([]model.TeamStatOfMatch, bool)
}

type gameOfMatchDao struct {
	db *sql.DB
}

func NewGameOfMatchDao(db *sql.DB) GameOfMatchDao {
	return &gameOfMatchDao{db: db}
}

const (
	ColumnWithoutIdOfGameOfMatch = "`name`, `home_team_id`, `away_team_id`, `club_id`, `ground_id`, `home_jersey_color`, `away_jersey_color`, `open_time`, `enroll_start_time`, `enroll_end_time`, `enroll_quota`, `enroll_count`, `rent_cost`, `duration`, `create_time`, `status`"
	ColumnOfGameOfMatch          = "`id`, " + ColumnWithoutIdOfGameOfMatch
	TableNameOfGameOfMatch       = "game_of_match"
)

//根据id获取比赛信息
func (c *gameOfMatchDao) GetById(id int) (*model.GameOfMatch, bool) {
	if gameOfMatchs, ok := c.queryGameOfMatch(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfGameOfMatch, TableNameOfGameOfMatch), id); ok && len(gameOfMatchs) == 1 {
		return &gameOfMatchs[0], true
	}
	return nil, false
}

func buildQueryGameOfMatchSql(returnColumn string, clubId sql.NullInt64, startOpenTime, endOpenTime sql.NullString, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfGameOfMatch))
	var args []interface{}
	if clubId.Valid {
		querySql.WriteString(" and club_id=?")
		args = append(args, clubId)
	}
	if startOpenTime.Valid {
		querySql.WriteString(" and open_time>=?")
		args = append(args, startOpenTime)
	}
	if endOpenTime.Valid {
		querySql.WriteString(" and open_time<=?")
		args = append(args, endOpenTime)
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
func (c *gameOfMatchDao) QueryGameOfMatch(clubId sql.NullInt64, startOpenTime, endOpenTime sql.NullString, status []string, limit, offset int) ([]model.GameOfMatch, bool) {
	querySql, args := buildQueryGameOfMatchSql(ColumnOfGameOfMatch, clubId, startOpenTime, endOpenTime, status)
	querySql += " order by create_time desc, id limit ? offset ?"
	args = append(args, limit, offset)
	return c.queryGameOfMatch(querySql, args...)
}

//搜索计数
func (c *gameOfMatchDao) QueryCount(clubId sql.NullInt64, startOpenTime, endOpenTime sql.NullString, status []string) int {
	querySql, args := buildQueryGameOfMatchSql("count(*)", clubId, startOpenTime, endOpenTime, status)
	stmt, err := c.db.Prepare(querySql)
	if err != nil {
		log.Printf("预编译%s.QueryCount语句出错，err: %s\n", TableNameOfGameOfMatch, err.Error())
		return 0
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("%s.QueryCount查询出错，err: %s\n", TableNameOfGameOfMatch, err.Error())
		return 0
	}
	if !rows.Next() {
		return 0
	}
	count := 0
	err = rows.Scan(&count)
	if err != nil {
		log.Printf("%s.QueryCount获取数据出错，err: %s\n", TableNameOfGameOfMatch, err.Error())
	}
	return count
}

func (c *gameOfMatchDao) queryGameOfMatch(query string, args ...interface{}) ([]model.GameOfMatch, bool) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Printf("预编译%s.queryGameOfMatch语句出错，err: %s\n", TableNameOfGameOfMatch, err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("%s.queryGameOfMatch查询出错，err: %s\n", TableNameOfGameOfMatch, err.Error())
		return nil, false
	}

	gameOfMatchs := make([]model.GameOfMatch, 0)
	for rows.Next() {
		gameOfMatch := model.GameOfMatch{}
		err = rows.Scan(&gameOfMatch.Id, &gameOfMatch.Name, &gameOfMatch.HomeTeamId, &gameOfMatch.AwayTeamId, &gameOfMatch.ClubId, &gameOfMatch.GroundId, &gameOfMatch.HomeJerseyColor, &gameOfMatch.AwayJerseyColor, &gameOfMatch.OpenTime, &gameOfMatch.EnrollStartTime, &gameOfMatch.EnrollEndTime, &gameOfMatch.EnrollQuota, &gameOfMatch.EnrollCount, &gameOfMatch.RentCost, &gameOfMatch.Duration, &gameOfMatch.CreateTime, &gameOfMatch.Status)
		if err != nil {
			log.Printf("%s.queryGameOfMatch获取数据出错，err: %s\n", TableNameOfGameOfMatch, err.Error())
			return nil, false
		}
		gameOfMatchs = append(gameOfMatchs, gameOfMatch)
	}

	return gameOfMatchs, true
}

func (c *gameOfMatchDao) Insert(gameOfMatch *model.GameOfMatch) (*model.GameOfMatch, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?)", TableNameOfGameOfMatch, ColumnWithoutIdOfGameOfMatch))
	if err != nil {
		log.Printf("预编译插入game_of_match语句出错，err: %s\n", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(gameOfMatch.Name, gameOfMatch.HomeTeamId, gameOfMatch.AwayTeamId, gameOfMatch.ClubId, gameOfMatch.GroundId, gameOfMatch.HomeJerseyColor, gameOfMatch.AwayJerseyColor, gameOfMatch.OpenTime, gameOfMatch.EnrollStartTime, gameOfMatch.EnrollEndTime, gameOfMatch.EnrollQuota, gameOfMatch.EnrollCount, gameOfMatch.RentCost, gameOfMatch.Duration, gameOfMatch.CreateTime, gameOfMatch.Status)
	if err != nil {
		log.Printf("插入game_of_match出错，err: %s\n", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入game_of_match.id出错，err: %s\n", err.Error())
		return nil, false
	}
	gameOfMatch.Id = lastInsertId
	return gameOfMatch, true
}

//更新比赛状态
func (c *gameOfMatchDao) UpdateStatus(id int, status, oldStatus string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `status`=? where id=? and `status`=?", TableNameOfGameOfMatch))
	if err != nil {
		log.Printf("预编译更新game_of_match.status语句出错，err: %s\n", err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(status, id, oldStatus)
	if err != nil {
		log.Printf("更新game_of_match.status出错，err: %s\n", err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新game_of_match.status影响行数出错，err: %s\n", err.Error())
		return 0, false
	}
	return rowsAffected, true
}

//更新已报名人数
func (c *gameOfMatchDao) UpdateEnrollCount(id, count int) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `enroll_count`=? where id=? and `enroll_quota`>=?", TableNameOfGameOfMatch))
	if err != nil {
		log.Printf("预编译更新game_of_match.enroll_count语句出错，err: %s\n", err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(count, id, count)
	if err != nil {
		log.Printf("更新game_of_match.enroll_count出错，err: %s\n", err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新game_of_match.enroll_count影响行数出错，err: %s\n", err.Error())
		return 0, false
	}
	return rowsAffected, true
}

const (
	ColumnWithoutIdOfEnrollOfMatch = "`match_id`, `player_id`, `temporary_player`, `create_time`, `status`"
	ColumnOfEnrollOfMatch          = "`id`, " + ColumnWithoutIdOfEnrollOfMatch
	TableNameOfEnrollOfMatch       = "enroll_of_match"
)

func (c *gameOfMatchDao) QueryEnrollOfMatch(matchId int, playerId sql.NullInt64, status []string) ([]model.EnrollOfMatch, bool) {
	querySql, args := buildQueryEnrollOfMatchSql(ColumnOfEnrollOfMatch, matchId, playerId, status)
	querySql += " order by create_time desc, id"
	return c.queryEnrollOfMatch(querySql, args...)
}

//保存球员比赛报名记录
func (c *gameOfMatchDao) InsertEnrollOfMatch(enroll *model.EnrollOfMatch) (*model.EnrollOfMatch, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?)", TableNameOfEnrollOfMatch, ColumnWithoutIdOfEnrollOfMatch))
	if err != nil {
		log.Printf("预编译插入%s语句出错，err: %s\n", TableNameOfEnrollOfMatch, err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(enroll.MatchId, enroll.PlayerId, enroll.TemporaryPlayer, enroll.CreateTime, enroll.Status)
	if err != nil {
		log.Printf("插入%s出错，err: %s\n", TableNameOfEnrollOfMatch, err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入%s.id出错，err: %s\n", TableNameOfEnrollOfMatch, err.Error())
		return nil, false
	}
	enroll.Id = int(lastInsertId)
	return enroll, true
}

//更新球员比赛报名状态
func (c *gameOfMatchDao) UpdateEnrollOfMatchStatus(id int64, status, oldStatus string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `status`=? where id=? and `status`=?", TableNameOfEnrollOfMatch))
	sqlMsg := fmt.Sprintf("%s.status", TableNameOfEnrollOfMatch)
	if err != nil {
		log.Printf("预编译更新%s语句出错，err: %s\n", sqlMsg, err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(status, id, oldStatus)
	if err != nil {
		log.Printf("更新%s出错，err: %s\n", sqlMsg, err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新%s影响行数出错，err: %s\n", sqlMsg, err.Error())
		return 0, false
	}
	return rowsAffected, true
}

func (c *gameOfMatchDao) queryEnrollOfMatch(query string, args ...interface{}) ([]model.EnrollOfMatch, bool) {
	stmt, err := c.db.Prepare(query)
	sqlMsg := fmt.Sprintf("%s.queryEnrollOfMatch", TableNameOfEnrollOfMatch)
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

	enrollOfMatchs := make([]model.EnrollOfMatch, 0)
	for rows.Next() {
		enrollOfMatch := model.EnrollOfMatch{}
		err = rows.Scan(&enrollOfMatch.Id, &enrollOfMatch.MatchId, &enrollOfMatch.PlayerId, &enrollOfMatch.TemporaryPlayer, &enrollOfMatch.CreateTime, &enrollOfMatch.Status)
		if err != nil {
			log.Printf("%s获取数据出错，err: %s\n", sqlMsg, err.Error())
			return nil, false
		}
		enrollOfMatchs = append(enrollOfMatchs, enrollOfMatch)
	}

	return enrollOfMatchs, true
}

func buildQueryEnrollOfMatchSql(returnColumn string, matchId int, playerId sql.NullInt64, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfEnrollOfMatch))
	var args []interface{}
	querySql.WriteString(" and match_id=?")
	args = append(args, matchId)
	if playerId.Valid {
		querySql.WriteString(" and player_id=?")
		args = append(args, playerId)
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

const (
	ColumnWithoutIdTeamStatOfMatch = "`match_id`, `type`, `team_id`, `score`, `rent_amount`, `shot`, `foul`, `free_kick`, `penalty_kick`, `offside`, `corner`, `pass`, `yellow_card`, `red_card`, `create_time`"
	ColumnTeamStatOfMatch          = "`id`, " + ColumnWithoutIdTeamStatOfMatch
	TableNameOfTeamStatOfMatch     = "team_stat_of_match"
	MatchTeamTypeHome              = "home"
	MatchTeamTypeAway              = "away"
)

func (c *gameOfMatchDao) InsertTeamStatOfMatch(stat *model.TeamStatOfMatch) (*model.TeamStatOfMatch, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?)", TableNameOfEnrollOfMatch, ColumnWithoutIdOfEnrollOfMatch))
	if err != nil {
		log.Printf("预编译插入%s语句出错，err: %s\n", TableNameOfTeamStatOfMatch, err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(stat.MatchId, stat.Ttype, stat.TeamId, stat.Score, stat.RentAmount, stat.Shot, stat.Foul, stat.FreeKick, stat.PenaltyKick, stat.Offside, stat.Corner, &stat.Pass, stat.YellowCard, stat.RedCard, stat.CreateTime)
	if err != nil {
		log.Printf("插入%s出错，err: %s\n", TableNameOfTeamStatOfMatch, err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入%s.id出错，err: %s\n", TableNameOfTeamStatOfMatch, err.Error())
		return nil, false
	}
	stat.Id = lastInsertId
	return stat, true
}

func (c *gameOfMatchDao) QueryTeamStatOfMatch(matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString, limit, offset int) ([]model.TeamStatOfMatch, bool) {
	querySql, args := buildQueryTeamStatOfMatchSql(ColumnTeamStatOfMatch, matchId, teamId, ttype, startTime, endTime)
	querySql += " order by create_time desc, id"
	return c.queryTeamStatOfMatch(querySql, args...)
}

func (c *gameOfMatchDao) queryTeamStatOfMatch(query string, args ...interface{}) ([]model.TeamStatOfMatch, bool) {
	stmt, err := c.db.Prepare(query)
	sqlMsg := fmt.Sprintf("%s.queryTeamOfMatch")
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

	teamStatOfMatchs := make([]model.TeamStatOfMatch, 0)
	for rows.Next() {
		teamStatOfMatch := model.TeamStatOfMatch{}
		err = rows.Scan(&teamStatOfMatch.Id, &teamStatOfMatch.MatchId, &teamStatOfMatch.Ttype, &teamStatOfMatch.TeamId, &teamStatOfMatch.Score, &teamStatOfMatch.RentAmount, &teamStatOfMatch.Shot, &teamStatOfMatch.Foul, &teamStatOfMatch.FreeKick, &teamStatOfMatch.PenaltyKick, &teamStatOfMatch.Offside, &teamStatOfMatch.Corner, &teamStatOfMatch.Pass, &teamStatOfMatch.YellowCard, &teamStatOfMatch.RedCard, &teamStatOfMatch.CreateTime)
		if err != nil {
			log.Printf("%s获取数据出错，err: %s\n", sqlMsg, err.Error())
			return nil, false
		}
		teamStatOfMatchs = append(teamStatOfMatchs, teamStatOfMatch)
	}

	return teamStatOfMatchs, true
}

func buildQueryTeamStatOfMatchSql(returnColumn string, matchId, teamId sql.NullInt64, ttype, startTime, endTime sql.NullString) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfTeamStatOfMatch))
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
