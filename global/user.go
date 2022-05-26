package global

const (
	CHATROOM = "room"
	RADIO    = "radio"
	ORIENT   = "orient"
)

type User struct {
	UID       uint64   `json:"uid"`
	PartnerID uint64   `json:"partner_id"`
	CompanyID uint64   `json:"company_id"`
	Name      string   `json:"name"`
	Follow    []uint64 `json:"follow"`
	Type      string   `json:"type"`
}

type UserData struct {
	AccessToken string `json:"access_token"`
	UID         uint64 `json:"uid"`
	PartnerID   uint64 `json:"partner_id"`
	CompanyID   uint64 `json:"company_id"`
	Name        string `json:"name"`
	CompanyIds  uint64 `json:"company_ids"`
	Follow      string `json:"follow"`
	Type        string `json:"type"`
	Data        []Data `json:"data"`
}
type Data struct {
	ID      uint64 `json:"id"`
	Content string `json:"content"`
}
