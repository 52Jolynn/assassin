package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
	"fmt"
	"strings"
)

type GroundDao interface {
	GetById(id int) (*model.Ground, bool)
	Insert(ground *model.Ground) (*model.Ground, bool)
	QueryGround(name, ttype sql.NullString, clubId sql.NullInt64, status []string, limit, offset int) ([]model.Ground, bool)
	QueryCount(name, ttype sql.NullString, clubId sql.NullInt64, status []string) int
	Update(ground *model.Ground) (int64, bool)
}

type groundDao struct {
	db *sql.DB
}

func NewGroundDao(db *sql.DB) GroundDao {
	return &groundDao{db: db}
}

const (
	ColumnWithoutIdOfGround = "`name`, `remark`, `type`, `club_id`, `create_time`, `status`"
	ColumnOfGround          = "`id`, " + ColumnWithoutIdOfGround
	TableNameOfGround       = "ground"
)

//根据id获取俱乐部场地信息
func (c *groundDao) GetById(id int) (*model.Ground, bool) {
	if grounds, ok := c.queryGround(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfGround, TableNameOfGround), id); ok && len(grounds) == 1 {
		return &grounds[0], true
	}
	return nil, false
}

func buildQueryGroundSql(returnColumn string, name, ttype sql.NullString, clubId sql.NullInt64, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfGround))
	var args []interface{}
	if name.Valid {
		querySql.WriteString(" and name=?")
		args = append(args, name)
	}
	if ttype.Valid {
		querySql.WriteString(" and `type`=?")
		args = append(args, ttype.String)
	}
	if clubId.Valid {
		querySql.WriteString(" and club_id=?")
		args = append(args, clubId.Int64)
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
func (c *groundDao) QueryGround(name, ttype sql.NullString, clubId sql.NullInt64, status []string, limit, offset int) ([]model.Ground, bool) {
	querySql, args := buildQueryGroundSql(ColumnOfGround, name, ttype, clubId, status)
	querySql += " order by create_time desc, id limit ? offset ?"
	args = append(args, limit, offset)
	return c.queryGround(querySql, args...)
}

//搜索计数
func (c *groundDao) QueryCount(name, ttype sql.NullString, clubId sql.NullInt64, status []string) int {
	querySql, args := buildQueryGroundSql("count(*)", name, ttype, clubId, status)
	stmt, err := c.db.Prepare(querySql)
	sqlErrMsg := fmt.Sprintf("%s.QueryCount", TableNameOfGround)
	if err != nil {
		log.Printf("预编译%s语句出错，err: %s\n", sqlErrMsg, err.Error())
		return 0
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("%s查询出错，err: %s\n", sqlErrMsg, err.Error())
		return 0
	}
	if !rows.Next() {
		return 0
	}
	count := 0
	err = rows.Scan(&count)
	if err != nil {
		log.Printf("%s获取数据出错，err: %s\n", sqlErrMsg, err.Error())
	}
	return count
}

func (c *groundDao) queryGround(query string, args ...interface{}) ([]model.Ground, bool) {
	stmt, err := c.db.Prepare(query)
	sqlErrMsg := fmt.Sprintf("%s.queryGround", TableNameOfGround)
	if err != nil {
		log.Printf("预编译%s语句出错，err: %s\n", sqlErrMsg, err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("%s查询出错，err: %s\n", sqlErrMsg, err.Error())
		return nil, false
	}

	grounds := make([]model.Ground, 0)
	for rows.Next() {
		ground := model.Ground{}
		err = rows.Scan(&ground.Id, &ground.Name, &ground.Remark, &ground.Ttype, &ground.ClubId, &ground.CreateTime, &ground.Status)
		if err != nil {
			log.Printf("%s获取数据出错，err: %s\n", sqlErrMsg, err.Error())
			return nil, false
		}
		grounds = append(grounds, ground)
	}

	return grounds, true
}

func (c *groundDao) Insert(ground *model.Ground) (*model.Ground, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?)", TableNameOfGround, ColumnWithoutIdOfGround))
	if err != nil {
		log.Printf("预编译插入%s语句出错，err: %s\n", TableNameOfGround, err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(ground.Name, ground.Remark, ground.Ttype, ground.ClubId, ground.CreateTime, ground.Status)
	if err != nil {
		log.Printf("插入%s出错，err: %s\n", TableNameOfGround, err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入%s.id出错，err: %s\n", TableNameOfGround, err.Error())
		return nil, false
	}
	ground.Id = int(lastInsertId)
	return ground, true
}

func (c *groundDao) Update(ground *model.Ground) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `name`=?, `remark`=?, `type`=?, `club_id`=?, `create_time`=?, `status`=? where id=?", TableNameOfGround))
	if err != nil {
		log.Printf("预编译更新%s语句出错，err: %s\n", TableNameOfGround, err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(ground.Name, ground.Remark, ground.Ttype, ground.ClubId, ground.CreateTime, ground.Status, ground.Id)
	if err != nil {
		log.Printf("更新%s出错，err: %s\n", TableNameOfGround, err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新%s影响行数出错，err: %s\n", TableNameOfGround, err.Error())
		return 0, false
	}
	return rowsAffected, true
}
