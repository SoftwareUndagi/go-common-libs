package coremodel

//tablenameCommonSequence nama table . di constant untuk optimasi
const tablenameCommonSequence = "ct_common_sequence"

//CommonSequence table: ct_common_sequence
type CommonSequence struct {
	//ID nama sequence, column: sequence_name
	ID string `gorm:"column:sequence_name;primary_key" json:"id"`
	//LatestSequence sequence terakhir dari data, column: latest_seq
	LatestSequence int64 `gorm:"column:latest_seq" json:"latestSequence"`
	//Remark catatan untuk sequence, column: seq_remark
	Remark string `gorm:"column:seq_remark" json:"remark"`
}

//TableName table name for struct CommonSequence
func (u CommonSequence) TableName() string {
	return tablenameCommonSequence
}
