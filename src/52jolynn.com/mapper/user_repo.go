package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
	"fmt"
	"strings"
)

type UserDao interface {
	GetById(id int) (*model.User, bool)
	GetByMobile(mobile string) (*model.User, bool)
	ExistsByMobile(mobile string) (ok, exists bool)
	Insert(user *model.User) (*model.User, bool)
	UpdatePasswd(uid int64, passwd string) (int64, bool)
	UpdateStatus(uid int64, status string) (int64, bool)
}

type userDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) UserDao {
	return &userDao{db: db}
}

const (
	ColumnWithoutIdOfUser = "`mobile`, `passwd`, `wx_open_id`, `nickname`, `create_time`, `last_active_time`, `status`"
	ColumnOfUser          = "`id`, " + ColumnWithoutIdOfUser
	TableNameOfUser       = "user"
)

//根据id获取用户信息
func (c *userDao) GetById(id int) (*model.User, bool) {
	if users, ok := c.queryUser(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfUser, TableNameOfUser), id); ok && len(users) == 1 {
		return &users[0], true
	}
	return nil, false
}

//根据mobile获取用户信息
func (c *userDao) GetByMobile(mobile string) (*model.User, bool) {
	if users, ok := c.queryUser(fmt.Sprintf("select %s from %s where `mobile` = ?", ColumnOfUser, TableNameOfUser), mobile); ok && len(users) == 1 {
		return &users[0], true
	}
	return nil, false
}

//根据mobile判断用户是否存在
func (c *userDao) ExistsByMobile(mobile string) (ok, exists bool) {
	if users, ok := c.queryUser(fmt.Sprintf("select %s from %s where `mobile` = ?", ColumnOfUser, TableNameOfUser), mobile); ok {
		return true, len(users) == 1
	}
	return false, false
}

func buildQueryUserSql(returnColumn string, mobile sql.NullString, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfUser))
	var args []interface{}
	if mobile.Valid {
		querySql.WriteString(" and mobile=?")
		args = append(args, mobile)
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

func (c *userDao) queryUser(query string, args ...interface{}) ([]model.User, bool) {
	stmt, err := c.db.Prepare(query)
	sqlErrMsg := fmt.Sprintf("%s.queryUser", TableNameOfUser)
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

	users := make([]model.User, 0)
	for rows.Next() {
		user := model.User{}
		err = rows.Scan(&user.Id, &user.Mobile, &user.Passwd, &user.WxOpenId, &user.NickName, &user.CreateTime, &user.LastActiveTime, &user.Status)
		if err != nil {
			log.Printf("%s获取数据出错，err: %s\n", sqlErrMsg, err.Error())
			return nil, false
		}
		users = append(users, user)
	}

	return users, true
}

func (c *userDao) Insert(user *model.User) (*model.User, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?)", TableNameOfUser, ColumnWithoutIdOfUser))
	if err != nil {
		log.Printf("预编译插入%s语句出错，err: %s\n", TableNameOfUser, err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.Mobile, user.Passwd, user.WxOpenId, user.NickName, user.CreateTime, user.Status)
	if err != nil {
		log.Printf("插入%s出错，err: %s\n", TableNameOfUser, err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入%s.id出错，err: %s\n", TableNameOfUser, err.Error())
		return nil, false
	}
	user.Id = lastInsertId
	return user, true
}

func (c *userDao) UpdatePasswd(uid int64, passwd string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `passwd`=? where id=?", TableNameOfUser))
	if err != nil {
		log.Printf("预编译更新%s密码语句出错，err: %s\n", TableNameOfUser, err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(uid, passwd)
	if err != nil {
		log.Printf("更新%s密码出错，err: %s\n", TableNameOfUser, err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新%s密码影响行数出错，err: %s\n", TableNameOfUser, err.Error())
		return 0, false
	}
	return rowsAffected, true
}

func (c *userDao) UpdateStatus(uid int64, status string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update `%s` set `status`=? where id=?", TableNameOfUser))
	if err != nil {
		log.Printf("预编译更新%s状态语句出错，err: %s\n", TableNameOfUser, err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(uid, status)
	if err != nil {
		log.Printf("更新%s状态出错，err: %s\n", TableNameOfUser, err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新%s状态影响行数出错，err: %s\n", TableNameOfUser, err.Error())
		return 0, false
	}
	return rowsAffected, true
}
