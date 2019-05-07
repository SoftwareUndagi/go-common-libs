package coremodel

import (
	"time"
)

//tablenameApplicationUser nama table . di constant untuk optimasi
const tablenameApplicationUser = "sec_user"

//ApplicationUser table: sec_user
type ApplicationUser struct {
	//ID primary key, column: pk
	ID int64 `gorm:"column:pk;AUTO_INCREMENT;primary_key" json:"id"`
	//RealName nama lengkap, column: real_name
	RealName string `gorm:"column:real_name" json:"realName"`
	//Username username untuk login, column: user_name
	Username string `gorm:"column:user_name" json:"username"`
	//DataStatus status data (A vs D), column: data_status
	DataStatus string `gorm:"column:data_status" json:"dataStatus"`
	//AvatarURL url avatar dari user, column: avatar_url
	AvatarURL string `gorm:"column:avatar_url" json:"avatarUrl"`
	//Password password user(kalau pakai metode password standard), column: password
	Password string `gorm:"column:password" json:"-"`
	//DefaultBranchCode kode cabang default(kode unit kerja - di abaikan kalau PNS), column: default_branch_code
	DefaultBranchCode string `gorm:"column:default_branch_code" json:"defaultBranchCode"`
	//Email email dari user, column: email
	Email string `gorm:"column:email" json:"email"`
	//EmployeeNo reference ke data pegawai, column: emp_no
	EmployeeNo string `gorm:"column:emp_no" json:"employeeNo"`
	//ExpiredDate waktu expired dari user, column: expired_date
	ExpiredDate time.Time `gorm:"column:expired_date" json:"expiredDate"`
	//FailedLoginCount count berapa kali gagal login.counter. kalau melebihi max, maka di lock, column: failed_login_attemps
	FailedLoginCount int32 `gorm:"column:failed_login_attemps" json:"failedLoginCount"`
	//LocaleCode locale code, column: locale
	LocaleCode string `gorm:"column:locale" json:"localeCode"`
	//LockedFlag flag di lock atau tidak, column: is_locked
	LockedFlag string `gorm:"column:is_locked" json:"lockedFlag"`
	//Phone1 phone1, column: phone1
	Phone1 string `gorm:"column:phone1" json:"phone1"`
	//Phone2 phone secondary, column: phone2
	Phone2 string `gorm:"column:phone2" json:"phone2"`
	//Remark catatan, column: remark
	Remark string `gorm:"column:remark" json:"remark"`
	//Uuid UUID, untuk integrasi misal dengan firebase, column: uuid
	UUID string `gorm:"column:uuid" json:"uuid"`
}

//TableName table name for struct ApplicationUser
func (u ApplicationUser) TableName() string {
	return tablenameApplicationUser
}
