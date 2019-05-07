package dao

import (
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/SoftwareUndagi/go-common-libs/common"
)

type baseModel struct {
	//CreatorName username creator data. column: creator_name
	CreatorName string `gorm:"column:creator_name" json:"creator"`
	//CreatedAt column createdAt
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	//CreatorIPAddress IP address of data creator
	CreatorIPAddress string `gorm:"column:creator_ip_address" json:"creatorIpAddress"`
	//UpdateAt column: updatedAt
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	//ModifiedBy column: modified_by
	ModifiedBy string `gorm:"column:modified_by" json:"modifiedBy"`
}

type sampleModel struct {
	//Code undefined , column: code
	Code string `gorm:"column:code;primary_key" json:"code"`
	//Description undefined , column: description
	Description string `gorm:"column:description" json:"description"`
	//TypeDescription undefined , column: typeDescription
	TypeDescription string `gorm:"column:typeDescription" json:"typeDescription"`
	//Icon undefined , column: icon
	Icon string `gorm:"column:icon" json:"icon"`

	//tes ref
	Refer SampleRefer `gorm:"foreignkey:Icon" json:"refer"`
	baseModel
}

type SampleRefer struct {
	//Code undefined , column: code
	Code string `gorm:"column:code;primary_key" json:"code2"`
	//Description undefined , column: description
	Description string `gorm:"column:description" json:"description2"`
}

//TableName table name for struct AccomodationType
func (u sampleModel) TableName() string {
	return "dodol"
}

func TestSample(t *testing.T) {
	common.CaptureLog(t).Release()
	logEntry := logrus.WithField("method", "TestSample")
	smpl := sampleModel{Code: "APART", Description: "dodol", baseModel: baseModel{CreatedAt: time.Now()}}
	analisa := analizeModel(reflect.TypeOf(smpl))
	logEntry.WithField("rslt", analisa).Infof("Selesai %s", analisa.modelType.Name())
}

func TestCheckPointerToArray(t *testing.T) {
	common.CaptureLog(t).Release()
	logEntry := logrus.WithField("method", "TestCheckPointerToArray")
	var smplRslt []sampleModel
	sampleAcceptArray(logEntry, &smplRslt)

}

func TestArray(t *testing.T) {
	common.CaptureLog(t).Release()
	logEntry := logrus.WithField("method", "TestArray")
	sampleWithInterface := sampleReturnArrayInterfaces()
	samplePlain := sampleReturnPlainArray()
	for _, simple := range *samplePlain {
		logEntry.WithField("data", simple).Infof("Code : %s", simple.Code)
	}
	logEntry.WithField("raw", sampleWithInterface).Infof("Ok selesai")

}

func sampleAcceptArray(logEntry *logrus.Entry, p *[]sampleModel) {
	println("Mulai")
	logEntry.Infof("dodol")

}

func sampleReturnArrayInterfaces() (rtvl interface{}) {
	rtvl = &[]sampleModel{sampleModel{Code: "DODOL"}, sampleModel{Code: "GARUT"}}
	return
}

func sampleReturnPlainArray() (rtvl *[]sampleModel) {
	rtvl = &[]sampleModel{sampleModel{Code: "DODOL"}, sampleModel{Code: "GARUT"}}
	return
}
