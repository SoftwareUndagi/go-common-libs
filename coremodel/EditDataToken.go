package coremodel

import "time"

//tablenameEditDataToken nama table . di constant untuk optimasi
const tablenameEditDataToken = "ct_edit_data_token"

//EditDataToken table: ct_edit_data_token
type EditDataToken struct {
	//Token token, column: id
	Token string `gorm:"column:id;primary_key" json:"token"`
	//UserName ID user owner data, column: user_name
	UserName string `gorm:"column:user_name" json:"userName"`
	//ActiveFlag flag data aktiv atau tidak, column: active_flag
	ActiveFlag string `gorm:"column:active_flag" json:"activeFlag"`
	//ObjectName object di edit(add edit delete), column: object_name
	ObjectName string `gorm:"column:object_name" json:"objectName"`
	//ObjectID ID object, as string, column: object_id_as_str
	ObjectID string `gorm:"column:object_id_as_str" json:"objectId"`
	//CreatedAt column : createdAt time when data was created
	CreatedAt *time.Time `gorm:"column:createdAt" json:"createdAt"`
	//UpdatedAt last update at column : updatedAt
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

//TableName table name for struct EditDataToken
func (u EditDataToken) TableName() string {
	return tablenameEditDataToken
}

//BeforeUpdate hook pre create/ update
func (u *EditDataToken) BeforeUpdate() (err error) {
	if u.CreatedAt == nil {
		tgl := time.Now()
		u.CreatedAt = &tgl
	}
	return
}
