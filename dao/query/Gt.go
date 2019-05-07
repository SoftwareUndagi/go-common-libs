package query

import "fmt"

//Gt query dengan operator >
type gt struct {
	//Field nama struct field atau nama column database untuk query
	Field string
	//Value value untuk query. ini di sesuaikan dengan definisi datatype sebaiknya
	Value interface{}
}

//GenerateSQL generate SQL statement untuk query >
func (p *gt) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	parameters = []interface{}{p.Value}
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	SQL = fmt.Sprintf("%s > ? ", dbCol)
	return
}

//FieldName untuk keseragaman dengan interface Q. nama field
func (p *gt) FieldName() (fieldOrOperatorName string) {
	return p.Field
}

//Gt generate interface for query >
func Gt(fieldName string, fieldValue interface{}) Q {
	return &gt{Field: fieldName, Value: fieldValue}
}
