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

const ColumnOfClub = "`id`, `name`, `remark`, `address`, `tel`, `create_time`, `status`"

//根据id获取俱乐部信息
func (c *clubDao) GetById(id int) (*model.Club, bool) {
	if clubs, ok := c.queryClub(fmt.Sprintf("select %s from club where `id` = ?", ColumnOfClub), id); ok && len(*clubs) == 1 {
		return &(*clubs)[0], true
	}
	return nil, false
}

//根据name获取俱乐部信息
func (c *clubDao) GetByName(name string) (*model.Club, bool) {
	if clubs, ok := c.queryClub(fmt.Sprintf("select %s from club where `name` = ?", ColumnOfClub), name); ok && len(*clubs) == 1 {
		return &(*clubs)[0], true
	}
	return nil, false
}

//根据name判断俱乐部是否存在
func (c *clubDao) ExistsByName(name string) (ok, exists bool) {
	if clubs, ok := c.queryClub(fmt.Sprintf("select %s from club where `name` = ?", ColumnOfClub), name); ok {
		return true, len(*clubs) == 1
	}
	return false, false
}

func buildQueryClubSql(returnColumn string, name *string, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from club where 1=1", returnColumn))
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
		log.Printf("预编译club.QueryCount语句出错，err: %s", err.Error())
		return 0
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("club.QueryCount查询出错，err: %s", err.Error())
		return 0
	}
	if !rows.Next() {
		return 0
	}
	count := 0
	err = rows.Scan(&count)
	if err != nil {
		log.Printf("club.QueryCount获取数据出错，err: %s", err.Error())
	}
	return count
}

func (c *clubDao) queryClub(query string, args ...interface{}) (*[]model.Club, bool) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Printf("预编译club.queryClub语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("club.queryClub查询出错，err: %s", err.Error())
		return nil, false
	}

	clubs := make([]model.Club, 0)
	for rows.Next() {
		club := model.Club{}
		err = rows.Scan(&club.Id, &club.Name, &club.Remark, &club.Address, &club.Tel, &club.CreateTime, &club.Status)
		if err != nil {
			log.Printf("club.queryClub获取数据出错，err: %s", err.Error())
			return nil, false
		}
		clubs = append(clubs, club)
	}

	return &clubs, true
}

func (c *clubDao) Insert(club *model.Club) (*model.Club, bool) {
	stmt, err := c.db.Prepare("insert into club (`name`, `remark`, `address`, `tel`, `create_time`, `status`) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("预编译插入club语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(club.Name, club.Remark, club.Address, club.Tel, club.CreateTime, club.Status)
	if err != nil {
		log.Printf("插入club出错，err: %s", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入id出错，err: %s", err.Error())
		return nil, false
	}
	club.Id = int(lastInsertId)
	return club, true
}

func (c *clubDao) Update(club *model.Club) (int64, bool)  {
	stmt, err := c.db.Prepare("update `club` set `name`=?, `remark`=?, `address`=?, `tel`=?, `create_time`=?, `status`=? where id=?")
	if err != nil {
		log.Printf("预编译更新club语句出错，err: %s", err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(club.Name, club.Remark, club.Address, club.Tel, club.CreateTime, club.Status, club.Id)
	if err != nil {
		log.Printf("更新club出错，err: %s", err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新club影响行数出错，err: %s", err.Error())
		return 0, false
	}
	return rowsAffected, true
}
