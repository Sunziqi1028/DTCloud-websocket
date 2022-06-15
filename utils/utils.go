package utils

import (
	"fmt"
	"gitee.com/ling-bin/netwebSocket/global"
	"os"
	"path/filepath"
	"strings"
	"time"

	"strconv"
)

// CheckUidUnique 校验UID 是否唯一
func CheckUidUnique(uid uint64) bool {
	if _, ok := global.GlobalUsers[uid]; ok {
		return false
	}
	return true
}

// CheckPartnerIDUnique 校验partner 是否唯一
func CheckPartnerIDUnique(partner_id uint64) bool {
	if _, ok := global.PartnerMap[partner_id]; ok {
		return false
	}
	return true
}

// ConvertString2IntSlice 字符转化成int 类型的slice
func ConvertString2IntSlice(s string) ([]uint64, error) {
	var intSlice []uint64
	fmt.Println(s)
	tmp := strings.Split(s, ",")
	for _, v := range tmp {
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		intSlice = append(intSlice, uint64(i))
	}
	return intSlice, nil
}

// ConvertTableName 转换表名
func ConvertTableName(beforeTableName string) (afterTableName string) {
	afterTableName = strings.Replace(beforeTableName, ".", "_", -1)
	return afterTableName
}

//ubuntu /opt/websocket/log/2020-05-28/01.log
//Windows c:/opt/websocket/log/2020-05-28/01.log

// MkdirLogDir  创建LOG path
func MkdirLogDir() (logPath string, err error) {
	path, _ := os.Getwd()
	nowDate := time.Now().Format("2006-01-02")
	logPath = filepath.Join(path, "log", nowDate)
	exist, err := IsExistFile(logPath)
	if !exist && err == nil {
		err = os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return logPath, nil
}

// CreateLogFile 创建Log文件
func CreateLogFile(logPath string) (*os.File, error) {
	nowHour := time.Now().Hour()
	nowHourStting := strconv.Itoa(nowHour)
	logFileName := logPath + "/" + nowHourStting + ".log"
	exist, err := IsExistFile(logFileName)
	if !exist && err == nil {
		f, err := os.Create(logFileName)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	return nil, err
}

// IsExistFile 判断文件是否存在
func IsExistFile(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// StoreDBMap 存储该DB对应所有的UID
func StoreDBMap(dbName string, uid uint64) {
	if oldUids, ok := global.OneDB2Uids[dbName]; ok {
		for _, oldUuid := range oldUids {
			if oldUuid == uid {
				continue
			} else {
				oldUids = append(oldUids, uid)
				global.OneDB2Uids[dbName] = oldUids
			}
		}
	} else {
		var newUids = []uint64{}
		newUids = append(newUids, uid)
		global.OneDB2Uids[dbName] = newUids
	}
}

// StoreCompanyIDMap 存储全部的用户信息
func StoreCompanyIDMap(companyId uint64, uid uint64) {
	if oldUids, ok := global.OneCompanyId2Uids[companyId]; ok {
		for _, oldUid := range oldUids {
			if oldUid == uid {
				continue
			}
		}
		oldUids = append(oldUids, uid)
		global.OneCompanyId2Uids[companyId] = oldUids
	} else {
		var newUids = []uint64{}
		newUids = append(newUids, uid)
		global.OneCompanyId2Uids[companyId] = newUids
	}
}

// StoreUuidOfUid 存储一个用户的多个UUID（窗口多开）
func StoreUuidOfUid(uid uint64, uuid uint64) {
	if oldUuids, ok := global.OneUId2Uuids[uid]; ok {
		for _, oldUuid := range oldUuids {
			if oldUuid == uuid {
				continue
			} else {
				oldUuids = append(oldUuids, uuid)
				global.OneUId2Uuids[uid] = oldUuids
				fmt.Println("utils.go , line:127, oneUid:", global.OneUId2Uuids)
			}
		}
	} else {
		var newUuids = []uint64{}
		newUuids = append(newUuids, uuid)
		global.OneUId2Uuids[uid] = newUuids
	}
}

// 用户离开，删除当前用户的所有信息
func DelGlobalUser(uid uint64) {
	delete(global.GlobalUsers, uid)
}

// DelUidOfDBMap 用户离开，删除当前DB map中的uid
func DelUidOfDBMap(uid uint64) {
	userInfo := global.GlobalUsers[uid]
	if userInfo != nil {
		uids := global.OneDB2Uids[userInfo.DatabaseSecret]
		var uidsNew = []uint64{}
		for k, uidTmp := range uids {
			if uid == uidTmp {
				uidsNew = append(uidsNew, uids[:k]...)
				uidsNew = append(uidsNew, uids[k+1:]...)
			}
		}
		global.OneDB2Uids[userInfo.DatabaseSecret] = uidsNew
	}
}

// DelUidOfCompanyID 用户离开，删除当前CompanyID 中的UID
func DelUidOfCompanyID(uid uint64) {
	userInfo := global.GlobalUsers[uid]
	if userInfo != nil {
		uids := global.OneCompanyId2Uids[userInfo.CompanyID]
		var uidsNew = []uint64{}
		for k, uidTmp := range uids {
			if uid == uidTmp {
				uidsNew = append(uidsNew, uids[:k]...)
				uidsNew = append(uidsNew, uids[k+1:]...)
			}
		}
		global.OneDB2Uids[userInfo.DatabaseSecret] = uidsNew
	}
}

// DelUuidOfUid 当前窗口用户离开，删除当前窗口的UUID
func DelUuidOfUid(uid uint64) {
	userInfo := global.GlobalUsers[uid]
	uuids := global.OneUId2Uuids[uid]
	var uuidsNew = []uint64{}
	for k, uuid := range uuids {
		if uuid == userInfo.UUID {
			uuidsNew = append(uuidsNew, uuids[:k]...)
			uuidsNew = append(uuidsNew, uuids[k+1:]...)
		}
	}
	global.OneUId2Uuids[uid] = uuidsNew
}
