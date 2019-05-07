package coremodel

//tablenameApplicationMenu nama table . di constant untuk optimasi
const tablenameApplicationMenu = "sec_menu"

//ApplicationMenu table: sec_menu
type ApplicationMenu struct {
	//ID id dari menu, column: pk
	ID int64 `gorm:"column:pk;AUTO_INCREMENT;primary_key" json:"id"`
	//Code kode menu, column: menu_code
	Code string `gorm:"column:menu_code" json:"code"`
	//HavePanelFlag flag apakah memiliki panel atau tiidak, column: is_have_panel
	HavePanelFlag string `gorm:"column:is_have_panel" json:"havePanelFlag"`
	//MenuCSS css menu(icon dsb), column: menu_css
	MenuCSS string `gorm:"column:menu_css" json:"menuCss"`
	//ParentID id induk dari menu, column: parent_id
	ParentID int64 `gorm:"column:parent_id" json:"parentId"`
	//Label label dari menu, column: menu_label
	Label string `gorm:"column:menu_label" json:"label"`
	//MenuTreeCode tree code dari menu, column: menu_tree_code
	MenuTreeCode string `gorm:"column:menu_tree_code" json:"menuTreeCode"`
	//OrderNumber urutan data, column: order_no
	OrderNumber int32 `gorm:"column:order_no" json:"orderNumber"`
	//I18nKey key internalization, column: i18n_key
	I18nKey string `gorm:"column:i18n_key" json:"i18nKey"`
	//RoutePath path dari handler, column: route_path
	RoutePath string `gorm:"column:route_path" json:"routePath"`
	//AdditionalParameter ,, column: additional_param
	AdditionalParameter string `gorm:"column:additional_param" json:"additionalParameter"`
	//StatusCode status data, column: data_status
	StatusCode string `gorm:"column:data_status" json:"statusCode"`
	//TreeLevelPosition level menu. pada level berapa data berada, column: tree_level_position
	TreeLevelPosition int32 `gorm:"column:tree_level_position" json:"treeLevelPosition"`
}

//TableName table name for struct ApplicationMenu
func (u ApplicationMenu) TableName() string {
	return tablenameApplicationMenu
}
