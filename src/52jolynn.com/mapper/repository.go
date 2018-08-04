package mapper

import (
	"52jolynn.com/model"
	"database/sql"
)

type ClubDao interface {
	Insert(club model.Club) (*model.Club, bool)
}

type clubDao struct {
	db *sql.DB
}

func NewClubDao(db *sql.DB) *clubDao {
	return &clubDao{db: db}
}

func (c *clubDao) Insert(club model.Club) (*model.Club, bool) {
	return nil, false
}
