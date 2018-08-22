package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
	"fmt"
	"strings"
)

type PlayerDao interface {
	GetById(id int) (*model.Player, bool)
	Insert(player *model.Player) (*model.Player, bool)
	QueryPlayer(username, name, mobile, level sql.NullString, age sql.NullInt64, status []string, limit, offset int) ([]model.Player, bool)
	QueryCount(username, name, mobile, level sql.NullString, age sql.NullInt64, status []string) int

	JoinTeam(pt *model.PlayerOfTeam) (*model.PlayerOfTeam, bool)
}

type playerDao struct {
	db *sql.DB
}

func NewPlayerDao(db *sql.DB) PlayerDao {
	return &playerDao{db: db}
}

const (
	ColumnWithoutIdOfPlayer = "`uid`, `name`, `remark`, `mobile`, `pos`, `height`, `age`, `pass_val`, `shot_val`, `strength_val`, `dribble_val`, `speed_val`, `tackle_val`, `head_val`, `throwing_val`, `reaction_val`, `create_time`, `status`, `level`"
	ColumnOfPlayer          = "`id`, " + ColumnWithoutIdOfPlayer
	TableNameOfPlayer       = "player"
)

//根据id获取球员信息
func (c *playerDao) GetById(id int) (*model.Player, bool) {
	if players, ok := c.queryPlayer(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfPlayer, TableNameOfPlayer), id); ok && len(players) == 1 {
		return &players[0], true
	}
	return nil, false
}

func buildQueryPlayerSql(returnColumn string, username, name, mobile, level sql.NullString, age sql.NullInt64, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfPlayer))
	var args []interface{}
	if username.Valid {
		querySql.WriteString(" and username=?")
		args = append(args, name)
	}
	if name.Valid {
		querySql.WriteString(" and name=?")
		args = append(args, name)
	}
	if mobile.Valid {
		querySql.WriteString(" and mobile=?")
		args = append(args, name)
	}
	if level.Valid {
		querySql.WriteString(" and level=?")
		args = append(args, name)
	}
	if age.Valid {
		querySql.WriteString(" and age=?")
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
func (c *playerDao) QueryPlayer(username, name, mobile, level sql.NullString, age sql.NullInt64, status []string, limit, offset int) ([]model.Player, bool) {
	querySql, args := buildQueryPlayerSql(ColumnOfPlayer, username, name, mobile, level, age, status)
	querySql += " order by create_time desc, id limit ? offset ?"
	args = append(args, limit, offset)
	return c.queryPlayer(querySql, args...)
}

//搜索计数
func (c *playerDao) QueryCount(username, name, mobile, level sql.NullString, age sql.NullInt64, status []string) int {
	querySql, args := buildQueryPlayerSql("count(*)", username, name, mobile, level, age, status)
	stmt, err := c.db.Prepare(querySql)
	if err != nil {
		log.Printf("预编译playerDao.QueryCount语句出错，err: %s", err.Error())
		return 0
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("playerDao.QueryCount查询出错，err: %s", err.Error())
		return 0
	}
	if !rows.Next() {
		return 0
	}
	count := 0
	err = rows.Scan(&count)
	if err != nil {
		log.Printf("playerDao.QueryCount获取数据出错，err: %s", err.Error())
	}
	return count
}

func (c *playerDao) queryPlayer(query string, args ...interface{}) ([]model.Player, bool) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Printf("预编译playerDao.queryPlayer语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("playerDao.queryPlayer查询出错，err: %s", err.Error())
		return nil, false
	}

	players := make([]model.Player, 0)
	for rows.Next() {
		player := model.Player{}
		err = rows.Scan(&player.Id, &player.Uid, &player.Name, &player.Remark, &player.Mobile, &player.Pos, &player.Height, &player.Age, &player.PassVal, &player.ShotVal, &player.StrengthVal, &player.DribbleVal, &player.SpeedVal, &player.TackleVal, &player.HeadVal, &player.ThrowingVal, &player.ReactionVal, &player.CreateTime, &player.Status, &player.Level)
		if err != nil {
			log.Printf("playerDao.queryPlayer获取数据出错，err: %s", err.Error())
			return nil, false
		}
		players = append(players, player)
	}

	return players, true
}

func (c *playerDao) Insert(player *model.Player) (*model.Player, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", TableNameOfPlayer, ColumnWithoutIdOfPlayer))
	if err != nil {
		log.Printf("预编译插入player语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(player.Uid, player.Name, player.Remark, player.Mobile, player.Pos, player.Height, player.Age, player.PassVal, player.ShotVal, player.StrengthVal, player.DribbleVal, player.SpeedVal, player.TackleVal, player.HeadVal, player.ThrowingVal, player.ReactionVal, player.CreateTime, player.Status, player.Level)
	if err != nil {
		log.Printf("插入player出错，err: %s", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入player.id出错，err: %s", err.Error())
		return nil, false
	}
	player.Id = int(lastInsertId)
	return player, true
}

const (
	ColumnWithoutIdOfPlayerOfTeam = "`player_id`, `team_id`, `no`, `present_balance`, `used_amount`, `join_time`, `create_time`"
	TableNameOfPlayerOfTeam       = "player_of_team"
)

func (c *playerDao) JoinTeam(pt *model.PlayerOfTeam) (*model.PlayerOfTeam, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?, ?)", TableNameOfPlayerOfTeam, ColumnWithoutIdOfPlayerOfTeam))
	if err != nil {
		log.Printf("预编译插入player_of_team语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(pt.PlayerId, pt.TeamId, pt.No, pt.PresentBalance, pt.UsedAmount, pt.JoinTime, pt.CreateTime)
	if err != nil {
		log.Printf("插入player_of_team出错，err: %s", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入player_of_team.id出错，err: %s", err.Error())
		return nil, false
	}
	pt.Id = int(lastInsertId)
	return pt, true
}
