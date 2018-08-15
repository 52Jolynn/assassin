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
	QueryGround(name, ttype *string, clubId *int, status []string, limit, offset int) (*[]model.Ground, bool)
	QueryCount(name, ttype *string, clubId *int, status []string) int
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
	if grounds, ok := c.queryGround(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfGround, TableNameOfGround), id); ok && len(*grounds) == 1 {
		return &(*grounds)[0], true
	}
	return nil, false
}

func buildQuerygroundSql(returnColumn string, name, ttype *string, clubId *int, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfGround))
	var args []interface{}
	if name != nil {
		querySql.WriteString(" and name=?")
		args = append(args, name)
	}
	if ttype != nil {
		querySql.WriteString(" and `type`=?")
		args = append(args, *ttype)
	}
	if clubId != nil {
		querySql.WriteString(" and club_id=?")
		args = append(args, *clubId)
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
func (c *groundDao) QueryGround(name, ttype *string, clubId *int, status []string, limit, offset int) (*[]model.Ground, bool) {
	querySql, args := buildQuerygroundSql(ColumnOfGround, name, ttype, clubId, status)
	querySql += " order by create_time desc, id limit ? offset ?"
	args = append(args, limit, offset)
	return c.queryGround(querySql, args...)
}

//搜索计数
func (c *groundDao) QueryCount(name, ttype *string, clubId *int, status []string) int {
	querySql, args := buildQuerygroundSql("count(*)", name, ttype, clubId, status)
	stmt, err := c.db.Prepare(querySql)
	if err != nil {
		log.Printf("预编译groundDao.QueryCount语句出错，err: %s", err.Error())
		return 0
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("groundDao.QueryCount查询出错，err: %s", err.Error())
		return 0
	}
	if !rows.Next() {
		return 0
	}
	count := 0
	err = rows.Scan(&count)
	if err != nil {
		log.Printf("groundDao.QueryCount获取数据出错，err: %s", err.Error())
	}
	return count
}

func (c *groundDao) queryGround(query string, args ...interface{}) (*[]model.Ground, bool) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Printf("预编译groundDao.queryGround语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("groundDao.queryGround查询出错，err: %s", err.Error())
		return nil, false
	}

	grounds := make([]model.Ground, 0)
	for rows.Next() {
		ground := model.Ground{}
		err = rows.Scan(&ground.Id, &ground.Name, &ground.Remark, &ground.Ttype, &ground.ClubId, &ground.CreateTime, &ground.Status)
		if err != nil {
			log.Printf("groundDao.queryGround获取数据出错，err: %s", err.Error())
			return nil, false
		}
		grounds = append(grounds, ground)
	}

	return &grounds, true
}

func (c *groundDao) Insert(ground *model.Ground) (*model.Ground, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?)", TableNameOfGround, ColumnWithoutIdOfGround))
	if err != nil {
		log.Printf("预编译插入ground语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(ground.Name, ground.Remark, ground.Ttype, ground.ClubId, ground.CreateTime, ground.Status)
	if err != nil {
		log.Printf("插入ground出错，err: %s", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入ground.id出错，err: %s", err.Error())
		return nil, false
	}
	ground.Id = int(lastInsertId)
	return ground, true
}

func (c *groundDao) Update(ground *model.Ground) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `name`=?, `remark`=?, `type`=?, `club_id`=?, `create_time`=?, `status`=? where id=?", TableNameOfGround))
	if err != nil {
		log.Printf("预编译更新ground语句出错，err: %s", err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(ground.Name, ground.Remark, ground.Ttype, ground.ClubId, ground.CreateTime, ground.Status, ground.Id)
	if err != nil {
		log.Printf("更新ground出错，err: %s", err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新ground影响行数出错，err: %s", err.Error())
		return 0, false
	}
	return rowsAffected, true
}
