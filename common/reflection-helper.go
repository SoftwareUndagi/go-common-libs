package common

import "reflect"

//GetReflectTypeOfStructObject read actual object type. this function will scan for actual object type  when wrapped as interface{}. actual object type could be determined
func GetReflectTypeOfStructObject(sampleModel interface{}) reflect.Type {
	//:= reflect.TypeOf(sampleModel)
	rv := reflect.ValueOf(sampleModel)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	return rv.Type()
}
