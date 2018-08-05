package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
)

type ClubDao interface {
	GetByName(name string) (*model.Club, bool)
	Insert(club *model.Club) (*model.Club, bool)
}

type clubDao struct {
	db *sql.DB
}

func NewClubDao(db *sql.DB) ClubDao {
	return &clubDao{db: db}
}

func (c *clubDao) GetByName(name string) (*model.Club, bool) {
	stmt, err := c.db.Prepare("select `id`, `name`, `remark`, `address`, `tel`, `create_time`, `status` from club where `name` = ?")
	if err != nil {
		log.Printf("预编译club.GetByName语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(name)
	if err != nil {
		log.Printf("club.GetByName查询出错，err: %s", err.Error())
		return nil, false
	}
	if !rows.Next() {
		return nil, false
	}

	club := model.Club{}
	err = rows.Scan(&club.Id, &club.Name, &club.Remark, &club.Address, &club.Tel, &club.CreateTime, &club.Status)
	if err != nil {
		log.Printf("club.GetByName获取数据出错，err: %s", err.Error())
		return nil, false
	}
	return &club, true
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
