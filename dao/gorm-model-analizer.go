package dao

import (
	"reflect"
	"strings"
)

//extractGormResult hasil pembacaan metadata
type extractGormResult struct {
	//columnName database column nanme
	columnName string
	//primaryKey primary key flag
	primaryKey bool
	//dataTypeName data type, ini string
	dataTypeName string
	//jsonName json serialization, if tag = - then zero lenth string
	jsonName string
	//name golang struct field name
	name string
}

//extractGormForeignKeyResult exract result for foreign key
type extractGormForeignKeyResult struct {
	//jsonName json serialization, if tag = - then zero lenth string
	jsonName string
	//name golang struct field name
	name string
	//destinationModelType tipe destination
	destinationModelType reflect.Type

	//referFieldName field name for refer
	referFieldName string
}

//analizeModelColumns membaca data columns + extract data primary key dari data
func analizeModelColumns(structType reflect.Type) (primaryKeyColumn extractGormResult, columns []extractGormResult, foreignKeys []extractGormForeignKeyResult) {
	for i := 0; i < structType.NumField(); i++ {
		theF := structType.Field(i)
		if theF.Type.Kind() == reflect.Struct {
			if !theF.Anonymous {
				foreignKey := parseForeignKey(theF)
				if len(foreignKey.referFieldName) > 0 {
					foreignKeys = append(foreignKeys, foreignKey)
				}
				continue // ini ada ke forein key berarti
				//println("Tes" + theF.Name)
			}
			lvl2Pk, lvl2Cols, foreignKeysLocal := analizeModelColumns(theF.Type)
			if len(lvl2Pk.columnName) > 0 {
				primaryKeyColumn = lvl2Pk
			}
			if len(lvl2Cols) > 0 {
				columns = append(columns, lvl2Cols...)
			}
			if len(foreignKeysLocal) > 0 {
				foreignKeys = append(foreignKeys, foreignKeysLocal...)
			}
			continue
		}
		rslt := parseFieldSimpleField(theF)
		if rslt.primaryKey {
			primaryKeyColumn = rslt
		}
		columns = append(columns, rslt)
	}
	return
}

//parseForeignKey parse foreign key dari struct
func parseForeignKey(theF reflect.StructField) (rtvl extractGormForeignKeyResult) {
	rtvl.name = theF.Name
	jsonTag := strings.TrimSpace(theF.Tag.Get("json"))
	if len(jsonTag) > 0 {
		rtvl.jsonName = jsonTag
	}
	gormTag := theF.Tag.Get("gorm")
	if strings.Contains(gormTag, "foreignkey") {
		partisi := strings.Split(gormTag, ";")
		for _, g := range partisi {
			if strings.Contains(g, "foreignkey") {
				partisi2 := strings.Split(g, ":")
				if len(partisi2) == 2 {
					rtvl.referFieldName = strings.TrimSpace(partisi2[1])
				}
			}
		}
	}
	return
}

//parseFieldSimpleField parse simple field of struct. for field that type is struct, should not be handled by this method
//
func parseFieldSimpleField(theF reflect.StructField) (rtvl extractGormResult) {
	gormTag := theF.Tag.Get("gorm")

	rtvl = extractGormResult{name: theF.Name, primaryKey: strings.Contains(gormTag, "primary_key"), columnName: theF.Name, dataTypeName: theF.Type.Name(), jsonName: theF.Name}
	if strings.Contains(gormTag, "column") {
		partisi := strings.Split(gormTag, ";")
		for _, g := range partisi {
			if strings.Contains(g, "column") {
				partisi2 := strings.Split(g, ":")
				if len(partisi2) == 2 {
					rtvl.columnName = strings.TrimSpace(partisi2[1])
				}
			}
		}
	}

	jsonTag := strings.TrimSpace(theF.Tag.Get("json"))
	if len(jsonTag) > 0 {
		rtvl.jsonName = jsonTag
	}
	return

}
