package coredata


//SimpleUserData simple user data
type SimpleUserData struct {
	//ID id table sec_user
	ID int32
	//Username username. ini berisi sama dengan UserUUID kalau misal anonymous
	Username string
	//RealName dari column real_name
	RealName string
	//Email dari column: email
	Email string
	//UserUUID dari column uuid
	UserUUID string
	//Phone kalau user dengan phone auth , ini akan di ambil dari phone1
	Phone string
}