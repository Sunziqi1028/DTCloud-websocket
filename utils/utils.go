package utils

import (
	"fmt"
	"gitee.com/ling-bin/netwebSocket/global"
	"strings"

	"strconv"
)

func CheckUidUnique(uid uint64) bool {
	if _, ok := global.GlobalUsers[uid]; ok {
		return false
	}
	return true
}

func CheckPartnerIDUnique(partner_id uint64) bool {
	if _, ok := global.PartnerMap[partner_id]; ok {
		return false
	}
	return true
}

func ConvertString2IntSlice(s string) ([]uint64, error) {
	var followInt []uint64
	fmt.Println(s)
	tmp := strings.Split(s, ",")
	for _, v := range tmp {
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		followInt = append(followInt, uint64(i))
	}
	return followInt, nil
}
