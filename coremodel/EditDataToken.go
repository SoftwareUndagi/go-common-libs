package coremodel

//tablenameEditDataToken nama table . di constant untuk optimasi
const tablenameEditDataToken = "ct_edit_data_token"

//EditDataToken table: ct_edit_data_token
type EditDataToken struct {
	//Token token, column: id
	Token string `gorm:"column:id;primary_key" json:"token"`
	//UserID ID user owner data, column: user_id
	UserID int32 `gorm:"column:user_id" json:"userId"`
	//ActiveFlag flag data aktiv atau tidak, column: active_flag
	ActiveFlag string `gorm:"column:active_flag" json:"activeFlag"`
	//ObjectName object di edit(add edit delete), column: object_name
	ObjectName string `gorm:"column:object_name" json:"objectName"`
	//ObjectID ID object, as string, column: object_id_as_str
	ObjectID string `gorm:"column:object_id_as_str" json:"objectId"`
}

//TableName table name for struct EditDataToken
func (u EditDataToken) TableName() string {
	return tablenameEditDataToken
}
