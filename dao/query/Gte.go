package query

import "fmt"

//gte query untuk operator >=
type gte struct {
	//Field nama struct field atau nama column database untuk query
	Field string
	//Value value untuk query. ini di sesuaikan dengan definisi datatype sebaiknya
	Value interface{}
}

//GenerateSQL generate sql untuk >=
func (p *gte) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	parameters = []interface{}{p.Value}
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	SQL = fmt.Sprintf("%s >= ? ", dbCol)
	return
}

//FieldName untuk keseragaman dengan interface Q. nama field
func (p *gte) FieldName() (fieldOrOperatorName string) {
	return p.Field
}

//Gte generate struct gte. to meet interface def
func Gte(fieldName string, fieldValue interface{}) Q {
	return &gte{Field: fieldName, Value: fieldValue}
}
