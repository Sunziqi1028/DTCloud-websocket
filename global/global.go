package global

// GlobalUsers 存储全部用户信息
var GlobalUsers = make(map[uint64]*User, 1024)

// PartnerMap 用来检验业务ID是否唯一
var PartnerMap = make(map[uint64]uint64)

// UsersOfHTTP  存储全部用户信息
var UsersOfHTTP = make(map[uint64]*UserData, 1024)
