package coremodel

import "github.com/jinzhu/gorm"

//LookupDetail Simple table lookup.for table : m_lookup_details
type LookupDetail struct {
	//ID column id
	ID int32 `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	//DetailCode column: detail_code kode detail
	DetailCode string `gorm:"column:detail_code;" json:"detailCode"`
	//LovID column: lov_id
	LovID string `gorm:"column:lov_id;" json:"lovId"`
	//Label label for lookup column: lov_label
	Label string `gorm:"column:lov_label;" json:"label"`
	//Value1 label for value 1. arbitary data 1
	Value1 string `gorm:"column:val_1;" json:"value11"`
	//Value2 label for value 2. arbitary data 2
	Value2 string `gorm:"column:val_2;" json:"value12"`
	//I18nKey key internalization for lookup
	I18nKey string `gorm:"column:i18n_key" json:"i18nKey"`
	//SequenceNo sort no for lookup
	SequenceNo int16 `gorm:"column:seq_no" json:"sequenceNo"`
}

//TableName table name for m_lookup_details
func (p *LookupDetail) TableName(db *gorm.DB) (name string) {
	return "m_lookup_details"
}
