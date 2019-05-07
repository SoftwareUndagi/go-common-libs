package coremodel

//tablenameApplicationGroup nama table . di constant untuk optimasi
const tablenameApplicationGroup = "sec_group"

//ApplicationGroup table: sec_group
type ApplicationGroup struct {
	//ID id dari group, column: pk
	ID int32 `gorm:"column:pk;AUTO_INCREMENT;primary_key" json:"id"`
	//Code kode group, column: group_code
	Code string `gorm:"column:group_code" json:"code"`
	//Name nama group, column: group_name
	Name string `gorm:"column:group_name" json:"name"`
	//Remark catatan dari group, column: group_remark
	Remark string `gorm:"column:group_remark" json:"remark"`
	//UsageCounter berapa data yang sudah merefer ini , column: count_counter
	UsageCounter int32 `gorm:"column:count_counter" json:"usageCounter"`
}

//TableName table name for struct ApplicationGroup
func (u ApplicationGroup) TableName() string {
	return tablenameApplicationGroup
}
