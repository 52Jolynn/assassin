package mapper

import (
	"52jolynn.com/model"
	"database/sql"
	"log"
	"fmt"
	"strings"
)

type CouponDao interface {
	GetById(id int) (*model.Coupon, bool)
	Insert(coupon *model.Coupon) (*model.Coupon, bool)
	QueryCoupon(clubId, teamId *int, status []string, limit, offset int) (*[]model.Coupon, bool)
	QueryCount(clubId, teamId *int, status []string) int
	UpdateUsedTime(id int, usedTime, status string) (int64, bool)
	UpdateStatus(id int, status, oldStatus string) (int64, bool)
}

type couponDao struct {
	db *sql.DB
}

func NewCouponDao(db *sql.DB) CouponDao {
	return &couponDao{db: db}
}

const (
	ColumnWithoutIdOfCoupon = "`club_id`, `team_id`, `amount`, `least_amount`, `effective_time`, `used_time`, `expire_time`, `create_time`, `status`"
	ColumnOfCoupon          = "`id`, " + ColumnWithoutIdOfCoupon
	TableNameOfCoupon       = "coupon"
)

//根据id获取球队信息
func (c *couponDao) GetById(id int) (*model.Coupon, bool) {
	if coupons, ok := c.queryCoupon(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfCoupon, TableNameOfCoupon), id); ok && len(*coupons) == 1 {
		return &(*coupons)[0], true
	}
	return nil, false
}

func buildQueryCouponSql(returnColumn string, clubId, teamId *int, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfCoupon))
	var args []interface{}
	if clubId != nil {
		querySql.WriteString(" and club_id=?")
		args = append(args, clubId)
	}
	if teamId != nil {
		querySql.WriteString(" and team_id=?")
		args = append(args, clubId)
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
func (c *couponDao) QueryCoupon(clubId, teamId *int, status []string, limit, offset int) (*[]model.Coupon, bool) {
	querySql, args := buildQueryCouponSql(ColumnOfCoupon, clubId, teamId, status)
	querySql += " order by create_time desc, id limit ? offset ?"
	args = append(args, limit, offset)
	return c.queryCoupon(querySql, args...)
}

//搜索计数
func (c *couponDao) QueryCount(clubId, teamId *int, status []string) int {
	querySql, args := buildQueryCouponSql("count(*)", clubId, teamId, status)
	stmt, err := c.db.Prepare(querySql)
	if err != nil {
		log.Printf("预编译couponDao.QueryCount语句出错，err: %s", err.Error())
		return 0
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("couponDao.QueryCount查询出错，err: %s", err.Error())
		return 0
	}
	if !rows.Next() {
		return 0
	}
	count := 0
	err = rows.Scan(&count)
	if err != nil {
		log.Printf("couponDao.QueryCount获取数据出错，err: %s", err.Error())
	}
	return count
}

func (c *couponDao) queryCoupon(query string, args ...interface{}) (*[]model.Coupon, bool) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Printf("预编译couponDao.queryCoupon语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("couponDao.queryCoupon查询出错，err: %s", err.Error())
		return nil, false
	}

	coupons := make([]model.Coupon, 0)
	for rows.Next() {
		coupon := model.Coupon{}
		err = rows.Scan(&coupon.Id, &coupon.ClubId, &coupon.TeamId, &coupon.Amount, &coupon.LeastAmount, &coupon.EffectiveTime, &coupon.UsedTime, &coupon.ExpireTime, &coupon.CreateTime, &coupon.Status)
		if err != nil {
			log.Printf("couponDao.queryCoupon获取数据出错，err: %s", err.Error())
			return nil, false
		}
		coupons = append(coupons, coupon)
	}

	return &coupons, true
}

func (c *couponDao) Insert(coupon *model.Coupon) (*model.Coupon, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?, ?, ?)", TableNameOfCoupon, ColumnWithoutIdOfCoupon))
	if err != nil {
		log.Printf("预编译插入coupon语句出错，err: %s", err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(coupon.ClubId, coupon.TeamId, coupon.Amount, coupon.LeastAmount, coupon.EffectiveTime, coupon.UsedTime, coupon.ExpireTime, coupon.CreateTime, coupon.Status)
	if err != nil {
		log.Printf("插入coupon出错，err: %s", err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入coupon.id出错，err: %s", err.Error())
		return nil, false
	}
	coupon.Id = int(lastInsertId)
	return coupon, true
}

func (c *couponDao) UpdateUsedTime(id int, usedTime, status string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update %s set `used_time`=?, `status`=? where id=?", TableNameOfCoupon))
	if err != nil {
		log.Printf("预编译更新coupon.usedTime语句出错，err: %s", err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(usedTime, status, id)
	if err != nil {
		log.Printf("更新coupon.usedTime出错，err: %s", err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新coupon.usedTime影响行数出错，err: %s", err.Error())
		return 0, false
	}
	return rowsAffected, true
}

func (c *couponDao) UpdateStatus(id int, status, oldStatus string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update %s set `status`=? where id=? and `status`=?", TableNameOfCoupon))
	if err != nil {
		log.Printf("预编译更新coupon.status语句出错，err: %s", err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(status, id, oldStatus)
	if err != nil {
		log.Printf("更新coupon.status出错，err: %s", err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新coupon.status影响行数出错，err: %s", err.Error())
		return 0, false
	}
	return rowsAffected, true
}
