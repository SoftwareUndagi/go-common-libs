package dao

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/SoftwareUndagi/go-common-libs/common"
	"github.com/sirupsen/logrus"
)

func TestNestedParse(t *testing.T) {
	common.CaptureLog(t).Release()
	logEntry := logrus.WithField("method", "TestArray")
	sampleJSON := `{
        "username": "gede.sutarsa" , 
        "address": { "$like" : "tabanan%"} , 
        "level2.code": {"$eq" : "101"}
	}`
	jsonMap := make(map[string]interface{})
	json.Unmarshal([]byte(sampleJSON), &jsonMap)

	nestedField := jsonMap["address"]
	simpleField := jsonMap["username"]
	switch v := nestedField.(type) {
	case string:
		fmt.Println(v)
	case int32, int64:
		fmt.Println(v)
	case map[string]interface{}:
		fmt.Println("Ini sub json:  ")
	default:
		fmt.Println("unknown")
	}
	logEntry.WithField("parsed", jsonMap).WithField("nestedField", nestedField).WithField("simpleField", simpleField).Info("selesai")

}
