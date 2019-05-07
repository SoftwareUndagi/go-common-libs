package coremodel

//tablenameGroupToMenuAssignment nama table . di constant untuk optimasi
const tablenameGroupToMenuAssignment = "sec_menu_assignment"

//ApplicationGroupToMenuAssignment table: sec_menu_assignment
type ApplicationGroupToMenuAssignment struct {
	//ID primary key data, column: pk
	ID int64 `gorm:"column:pk;AUTO_INCREMENT;primary_key" json:"id"`
	//MenuID id dari menu, column: menu_id
	MenuID int64 `gorm:"column:menu_id" json:"menuId"`
	//GroupID id group, column: group_id
	GroupID int64 `gorm:"column:group_id" json:"groupId"`
	//AllowCreateFlag flag : ijinkan new data, column: is_allow_create
	AllowCreateFlag string `gorm:"column:is_allow_create" json:"allowCreateFlag"`
	//AllowEditFlag flag : ijinkan edit data, column: is_allow_edit
	AllowEditFlag string `gorm:"column:is_allow_edit" json:"allowEditFlag"`
	//AllowEraseFlag flag : ijinkan hapus data, column: is_allow_erase
	AllowEraseFlag string `gorm:"column:is_allow_erase" json:"allowEraseFlag"`
}

//TableName table name for struct GroupToMenuAssignment
func (u ApplicationGroupToMenuAssignment) TableName() string {
	return tablenameGroupToMenuAssignment
}
