/**
 * @Author: Sun
 * @Description:
 * @File:  db
 * @Version: 1.0.0
 * @Date: 2022/5/29 23:33
 */

package db

import (
	"fmt"
	"gitee.com/ling-bin/netwebSocket/global"
	"github.com/jinzhu/gorm"
	"log"
)

func NewPostgresql() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "host=%s  user=%s dbname=%s password=%s port=%s sslmode=disable",
		global.GlobalConfig.Host,
		global.GlobalConfig.Port,
		global.GlobalConfig.User,
		global.GlobalConfig.DataBase,
		global.GlobalConfig.Pass)
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		fmt.Println("connect postgres err:", err)
		return nil, err
	}

	return db, nil
}

type DBEngine struct {
	DBengine *gorm.DB
}

func (db *DBEngine) Insert(table string, data string) (err error) {
	err = db.DBengine.Table(table).Create(&data).Error // INSERT INTO res_partner(date,city, color) VALUES ('2022-05-26 01:06:04','上海',1)
	if err != nil {
		fmt.Println("create date err:", err)
		return err
	}
	return nil
}

func (db *DBEngine) Update(table string, ids []uint64, field string) (err error) {
	err = db.DBengine.Debug().Table(table).Where("id in ?", ids).
		Update(map[string]interface{}{"city": "?"}, field).Error // UPDATE res_partner SET date='2022-05-26 01:06:04', city='上海',color=3 WHERE id in (1,2)

	if err != nil {
		fmt.Println("update value err:", err)
		return err
	}
	return nil
}
