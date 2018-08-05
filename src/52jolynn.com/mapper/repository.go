package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
	"fmt"
)

type ClubDao interface {
	GetById(id int) (*model.Club, bool)
	GetByName(name string) (*model.Club, bool)
	ExistsByName(name string) (ok, exists bool)
	Insert(club *model.Club) (*model.Club, bool)
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

	var clubs []model.Club
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
