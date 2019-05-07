package coremodel

//tablenameUserToGroupAssignment nama table . di constant untuk optimasi
const tablenameUserToGroupAssignment = "sec_group_assignment"

//ApplicationUserToGroupAssignment table: sec_group_assignment
type ApplicationUserToGroupAssignment struct {
	//ID id data(surrogate key), column: pk
	ID int64 `gorm:"column:pk;AUTO_INCREMENT;primary_key" json:"id"`
	//GroupID id dari group, column: group_id
	GroupID int64 `gorm:"column:group_id" json:"groupId"`
	//UserID id dari user, column: user_id
	UserID int64 `gorm:"column:user_id" json:"userId"`
	//User refer with column user_id
	User ApplicationUser `gorm:"foreignkey:UserID" json:"user"`
	//Group refer with column group_id
	Group ApplicationGroup `gorm:"foreignkey:GroupID" json:"group"`
	//Creator column: creator_name
	Creator string `gorm:"column:creator_name" json:"creator"`
	//CreatorIpAddress column: creator_ip_address
	CreatorIPAddress string `gorm:"column:creator_ip_address" json:"creatorIpAddress"`
}

//TableName table name for struct UserToGroupAssignment
func (u ApplicationUserToGroupAssignment) TableName() string {
	return tablenameUserToGroupAssignment
}
