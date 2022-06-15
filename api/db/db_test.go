/**
 * @Author: Sun
 * @Description:
 * @File:  db_test
 * @Version: 1.0.0
 * @Date: 2022/6/10 11:18
 */

package db

import (
	"fmt"
	"gitee.com/ling-bin/netwebSocket/pkg"
	"testing"
	"time"
)

type ResUser struct {
	ID           int64     `gorm:"id"`
	Active       bool      `gorm:"active"`
	Login        string    `gorm:"login"`
	Password     string    `gorm:"password"`
	CompanyId    uint64    `gorm:"company_id"`
	PartnerId    uint64    `gorm:"partner_id"`
	CreateDate   time.Time `gorm:"create_date"`
	Signature    string    `gorm:"signature"`
	ActionId     uint64    `gorm:"action_id"`
	Share        bool      `gorm:"share"`
	DepartmentId uint64    `gorm:"department_id"`
	CreareUid    uint64    `gorm:"create_uid"`
	WriteUId     uint64    `gorm:"write_uid"`
	WriteDate    time.Time `gorm:"write_date"`
	TotpSecret   string    `gorm:"totp_secret"`
}

func TestDB(t *testing.T) {
	pkg.NewConf()
	db, err := NewPostgresql()
	if err != nil {
		t.Error("err", err)
	}
	var user ResUser
	db.Debug().Table("res_users").Where("login = public").First(&user)
	fmt.Println(user)
}
