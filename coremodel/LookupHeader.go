package coremodel

import (
	"github.com/jinzhu/gorm"
)

//LookupHeader table m_lookup_header
type LookupHeader struct {
	//ID code/ id of lookup column : lov_id
	ID string `gorm:"column:lov_id;primary_key" json:"id" `
	//Reamark remark for lookup. for maintence purpose. column lov_remark
	Remark string `gorm:"column:lov_remark" json:"remark"`
	//FlagCacheable Y or N. flag if lookup is cacheable on client
	FlagCacheable string `gorm:"column:is_cacheable" json:"flagCachable"`
	//FlagUseCustomSQL Y and N flag. is lookup using custom sql or from lookup detail
	FlagUseCustomSQL string `gorm:"column:is_use_custom_sql" json:"flagUseCustomSql"`
	//Version version of lookup. to force client reload if cache is expired
	Version string `gorm:"column:lov_version" json:"version"`
	//SQLForData sql for check LOV version. lookup version based on query result
	SQLForData *string `gorm:"column:sql_data" json:"sqlForData"`
	//SQLForDataFiltered sql for data filtered(id is passed with prev data)
	SQLForDataFiltered string `gorm:"column:sql_data_filtered" json:"sqlForDataFiltered"`
	//SQLForVersion sql for version of lookup
	SQLForVersion string `gorm:"column:sql_version" json:"sqlForVersion"`
}

//TableName table name for m_lookup_header
func (p *LookupHeader) TableName(db *gorm.DB) (name string) {
	return "m_lookup_header"
}

//BeforeUpdate hook pre create/ update
func (p *LookupHeader) BeforeUpdate() (err error) {
	if len(p.FlagUseCustomSQL) == 0 {
		if len(*p.SQLForData) > 0 {
			p.FlagUseCustomSQL = "Y"
		} else {
			p.FlagUseCustomSQL = "N"
		}

	}
	if len(p.FlagCacheable) == 0 {
		p.FlagCacheable = "Y"
	}
	if len(p.Version) == 0 {
		p.Version = "001"
	}
	return
}
