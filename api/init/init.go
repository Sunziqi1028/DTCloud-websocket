/**
 * @Author: Sun
 * @Description:
 * @File:  init
 * @Version: 1.0.0
 * @Date: 2022/5/29 23:19
 */

package init

import (
	"gitee.com/ling-bin/netwebSocket/api/db"
	"gitee.com/ling-bin/netwebSocket/global"
	"log"
)

func init() {
	err := setupDBEngine()
	if err != nil {
		log.Fatal("init.go line:28 connection postgresql err:", err)
	}
}

func setupDBEngine() error {
	var err error
	global.PostgreSqlDBEngine, err = db.NewPostgresql()
	if err != nil {
		return err
	}
	return nil
}
