package global

type CheckData struct {
	CheckValue  Value  `json:"value"`
	CheckParams Params `json:"params"`
}
type FundBillList struct {
	FundChannel string  `json:"fund_channel"`
	Amount      int     `json:"amount"`
	RealAmount  float64 `json:"real_amount"`
}
type VoucherDetailList struct {
	ID                         string  `json:"id"`
	PurchaseBuyerContribute    float64 `json:"purchase_buyer_contribute"`
	PurchaseMerchantContribute float64 `json:"purchase_merchant_contribute"`
	PurchaseAntContribute      float64 `json:"purchase_ant_contribute"`
}
type Value struct {
	Code               string              `json:"code"`
	Msg                string              `json:"msg"`
	TradeNo            string              `json:"trade_no"`
	GmtPayment         string              `json:"gmt_payment"`
	FundBillLists      []FundBillList      `json:"fund_bill_list"`
	StoreName          string              `json:"store_name"`
	BuyerUserID        string              `json:"buyer_user_id"`
	AsyncPaymentMode   string              `json:"async_payment_mode"`
	VoucherDetailLists []VoucherDetailList `json:"voucher_detail_list"`
}
type Context struct {
	Lang  string `json:"lang"`
	Email string `json:"email"`
}
type Params struct {
	Version     string    `json:"version"`      // 默认1.0
	Timestamp   string    `json:"timestamp"`    // 发送请求的时间，格式"yyyy-MM-dd HH:mm:ss
	Conn        string    `json:"conn"`         // 数据库连接
	AccessToken string    `json:"access_token"` // 用户授权
	URL         string    `json:"url"`          // 地址
	NotifyURL   string    `json:"notify_url"`   // 消息执行完后通知地址
	Dbname      string    `json:"dbname"`       // 执行的数据库
	UID         string    `json:"uid"`          // 用户ID
	PartnerID   string    `json:"partner_id"`   // partnerID
	CompanyID   string    `json:"company_id"`   // 公司ID
	Contexts    []Context `json:"context"`      // 上下文参数
	Method      string    `json:"method"`       // 接口名称 write,unlink,create
	Model       string    `json:"model"`        // 表名
	Field       string    `json:"field"`        // 操作的字段
	Body        string    `json:"body"`         // 描述
	Ids         string    `json:"ids"`          // 执行ID
	Fun         string    `json:"fun"`          // 执行的方法
	Type        string    `json:"type"`         // SQL(本地执行SQL),XMLRPC 发给远程地址 url 地址返回 T,log 写入本地文件 log 每小时一个文件名 年，月，日，时
}
