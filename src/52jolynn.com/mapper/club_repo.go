package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
	"fmt"
	"strings"
)

type ClubDao interface {
	GetById(id int) (*model.Club, bool)
	GetByName(name string) (*model.Club, bool)
	ExistsByName(name string) (ok, exists bool)
	Insert(club *model.Club) (*model.Club, bool)
	QueryClub(name *string, status []string, limit, offset int) (*[]model.Club, bool)
	QueryCount(name *string, status []string) int
	Update(club *model.Club) (int64, bool)
}

type clubDao struct {
	db *sql.DB
}

func NewClubDao(db *sql.DB) ClubDao {
	return &clubDao{db: db}
}

const (
	ColumnWithoutIdOfClub = "`name`, `remark`, `address`, `tel`, `create_time`, `status`"
	ColumnOfClub          = "`id`, " + ColumnWithoutIdOfClub
	TableNameOfClub       = "club"
)

//根据id获取俱乐部信息
func (c *clubDao) GetById(id int) (*model.Club, bool) {
	if clubs, ok := c.queryClub(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfClub, TableNameOfClub), id); ok && len(*clubs) == 1 {
		return &(*clubs)[0], true
	}
	return nil, false
}

//根据name获取俱乐部信息
func (c *clubDao) GetByName(name string) (*model.Club, bool) {
	if clubs, ok := c.queryClub(fmt.Sprintf("select %s from %s where `name` = ?", ColumnOfClub, TableNameOfClub), name); ok && len(*clubs) == 1 {
		return &(*clubs)[0], true
	}
	return nil, false
}

//根据name判断俱乐部是否存在
func (c *clubDao) ExistsByName(name string) (ok, exists bool) {
	if clubs, ok := c.queryClub(fmt.Sprintf("select %s from %s where `name` = ?", ColumnOfClub, TableNameOfClub), name); ok {
		return true, len(*clubs) == 1
	}
	return false, false
}

func buildQueryClubSql(returnColumn string, name *string, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfClub))
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
func (c *clubDao) QueryClub(name *string, status []string, limit, offset int) (*[]model.Club, bool) {
	querySql, args := buildQueryClubSql(ColumnOfClub, name, status)
	querySql += " order by create_time desc, id limit ? offset ?"
	args = append(args, limit, offset)
	return c.queryClub(querySql, args...)
}

//搜索计数
func (c *clubDao) QueryCount(name *string, status []string) int {
	querySql, args := buildQueryClubSql("count(*)", name, status)
	stmt, err := c.db.Prepare(querySql)
	if err != nil {
		log.Printf("预编译%s.QueryCount语句出错，err: %s\n", TableNameOfClub, err.Error())
		return 0
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("%s.QueryCount查询出错，err: %s\n", TableNameOfClub, err.Error())
		return 0
	}
	if !rows.Next() {
		return 0
	}
	count := 0
	err = rows.Scan(&count)
	if err != nil {
		log.Printf("%s.QueryCount获取数据出错，err: %s\n", TableNameOfClub, err.Error())
	}
	return count
}

func (c *clubDao) queryClub(query string, args ...interface{}) (*[]model.Club, bool) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Printf("预编译%s.queryClub语句出错，err: %s\n", TableNameOfClub, err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("%s.queryClub查询出错，err: %s\n", TableNameOfClub, err.Error())
		return nil, false
	}

	clubs := make([]model.Club, 0)
	for rows.Next() {
		club := model.Club{}
		err = rows.Scan(&club.Id, &club.Name, &club.Remark, &club.Address, &club.Tel, &club.CreateTime, &club.Status)
		if err != nil {
			log.Printf("%s.queryClub获取数据出错，err: %s\n", TableNameOfClub, err.Error())
			return nil, false
		}
		clubs = append(clubs, club)
	}

	return &clubs, true
}

func (c *clubDao) Insert(club *model.Club) (*model.Club, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?)", TableNameOfClub, ColumnWithoutIdOfClub))
	if err != nil {
		log.Printf("预编译插入%s语句出错，err: %s\n", TableNameOfClub, err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(club.Name, club.Remark, club.Address, club.Tel, club.CreateTime, club.Status)
	if err != nil {
		log.Printf("插入%s出错，err: %s\n", TableNameOfClub, err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入%s.id出错，err: %s\n", TableNameOfClub, err.Error())
		return nil, false
	}
	club.Id = int(lastInsertId)
	return club, true
}

func (c *clubDao) Update(club *model.Club) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `name`=?, `remark`=?, `address`=?, `tel`=?, `create_time`=?, `status`=? where id=?", TableNameOfClub))
	if err != nil {
		log.Printf("预编译更新%s语句出错，err: %s\n", TableNameOfClub, err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(club.Name, club.Remark, club.Address, club.Tel, club.CreateTime, club.Status, club.Id)
	if err != nil {
		log.Printf("更新%s出错，err: %s\n", TableNameOfClub, err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新%s影响行数出错，err: %s\n", TableNameOfClub, err.Error())
		return 0, false
	}
	return rowsAffected, true
}
