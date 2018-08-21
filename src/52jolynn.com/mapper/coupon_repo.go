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
	QueryCoupon(clubId, teamId sql.NullInt64, status []string, limit, offset int) ([]model.Coupon, bool)
	QueryCount(clubId, teamId sql.NullInt64, status []string) int
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

//根据id获取优惠信息
func (c *couponDao) GetById(id int) (*model.Coupon, bool) {
	if coupons, ok := c.queryCoupon(fmt.Sprintf("select %s from %s where `id` = ?", ColumnOfCoupon, TableNameOfCoupon), id); ok && len(coupons) == 1 {
		return &coupons[0], true
	}
	return nil, false
}

func buildQueryCouponSql(returnColumn string, clubId, teamId sql.NullInt64, status []string) (string, []interface{}) {
	querySql := strings.Builder{}
	querySql.WriteString(fmt.Sprintf("select %s from %s where 1=1", returnColumn, TableNameOfCoupon))
	var args []interface{}
	if clubId.Valid {
		querySql.WriteString(" and club_id=?")
		args = append(args, clubId)
	}
	if teamId.Valid {
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
func (c *couponDao) QueryCoupon(clubId, teamId sql.NullInt64, status []string, limit, offset int) ([]model.Coupon, bool) {
	querySql, args := buildQueryCouponSql(ColumnOfCoupon, clubId, teamId, status)
	querySql += " order by create_time desc, id limit ? offset ?"
	args = append(args, limit, offset)
	return c.queryCoupon(querySql, args...)
}

//搜索计数
func (c *couponDao) QueryCount(clubId, teamId sql.NullInt64, status []string) int {
	querySql, args := buildQueryCouponSql("count(*)", clubId, teamId, status)
	stmt, err := c.db.Prepare(querySql)
	sqlErrMsg := fmt.Sprintf("%s.QueryCount", TableNameOfCoupon)
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

func (c *couponDao) queryCoupon(query string, args ...interface{}) ([]model.Coupon, bool) {
	stmt, err := c.db.Prepare(query)
	sqlErrMsg := fmt.Sprintf("%s.queryCoupon", TableNameOfCoupon)
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

	coupons := make([]model.Coupon, 0)
	for rows.Next() {
		coupon := model.Coupon{}
		err = rows.Scan(&coupon.Id, &coupon.ClubId, &coupon.TeamId, &coupon.Amount, &coupon.LeastAmount, &coupon.EffectiveTime, &coupon.UsedTime, &coupon.ExpireTime, &coupon.CreateTime, &coupon.Status)
		if err != nil {
			log.Printf("%s获取数据出错，err: %s\n", sqlErrMsg, err.Error())
			return nil, false
		}
		coupons = append(coupons, coupon)
	}

	return coupons, true
}

func (c *couponDao) Insert(coupon *model.Coupon) (*model.Coupon, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("insert into %s (%s) values(?, ?, ?, ?, ?, ?, ?, ?, ?)", TableNameOfCoupon, ColumnWithoutIdOfCoupon))
	if err != nil {
		log.Printf("预编译插入%s语句出错，err: %s\n", TableNameOfCoupon, err.Error())
		return nil, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(coupon.ClubId, coupon.TeamId, coupon.Amount, coupon.LeastAmount, coupon.EffectiveTime, coupon.UsedTime, coupon.ExpireTime, coupon.CreateTime, coupon.Status)
	if err != nil {
		log.Printf("插入%s出错，err: %s\n", TableNameOfCoupon, err.Error())
		return nil, false
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入%s.id出错，err: %s\n", TableNameOfCoupon, err.Error())
		return nil, false
	}
	coupon.Id = int(lastInsertId)
	return coupon, true
}

func (c *couponDao) UpdateUsedTime(id int, usedTime, status string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update %s set `used_time`=?, `status`=? where id=?", TableNameOfCoupon))
	sqlErrMsg := fmt.Sprintf("%s.usedTime", TableNameOfCoupon)
	if err != nil {
		log.Printf("预编译更新%s语句出错，err: %s\n", sqlErrMsg, err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(usedTime, status, id)
	if err != nil {
		log.Printf("更新%s出错，err: %s\n", sqlErrMsg, err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新%s影响行数出错，err: %s\n", sqlErrMsg, err.Error())
		return 0, false
	}
	return rowsAffected, true
}

func (c *couponDao) UpdateStatus(id int, status, oldStatus string) (int64, bool) {
	stmt, err := c.db.Prepare(fmt.Sprintf("update %s set `status`=? where id=? and `status`=?", TableNameOfCoupon))
	sqlErrMsg := fmt.Sprintf("%s.status", TableNameOfCoupon)
	if err != nil {
		log.Printf("预编译更新%s语句出错，err: %s\n", sqlErrMsg, err.Error())
		return 0, false
	}
	defer stmt.Close()
	result, err := stmt.Exec(status, id, oldStatus)
	if err != nil {
		log.Printf("更新%s出错，err: %s\n", sqlErrMsg, err.Error())
		return 0, false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新%s影响行数出错，err: %s\n", sqlErrMsg, err.Error())
		return 0, false
	}
	return rowsAffected, true
}
