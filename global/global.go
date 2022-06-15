package global

import (
	"gitee.com/ling-bin/netwebSocket/pkg"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	SQL     = "sql"
	LOG     = "log"
	URL     = "url"
	XML_RPC = "xmlrpc"
	WRITE   = "write"
	CREATE  = "create"
	UNLINK  = "unlink"
)

var GlobalUsers = make(map[uint64]*User, 1024) // GlobalUsers 存储全部用户信息

var PartnerMap = make(map[uint64]uint64) // PartnerMap 用来检验业务ID是否唯一

var UsersOfHTTP = make(map[uint64]*UserData, 1024) // UsersOfHTTP  存储全部用户信息

var PostgreSqlDBEngine *gorm.DB // 全局gorm的客户端连接

var GlobalConfig = pkg.NewConf() // 全局配置文件

var OneCompanyId2Uids = make(map[uint64][]uint64) // 一个company_id对应该company_id所有的UID

var OneUId2Uuids = make(map[uint64][]uint64) // 一个UID对应多个UUIDS

var OneDB2Uids = make(map[string][]uint64) // 一个db对应该DB的所有UID
