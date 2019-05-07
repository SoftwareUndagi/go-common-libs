package query

import "fmt"

//lte query dengan operator <=
type lte struct {
	//Field nama struct field atau nama column database untuk query
	Field string
	//Value value untuk query. ini di sesuaikan dengan definisi datatype sebaiknya
	Value interface{}
}

//GenerateSQL generate sql untuk <=
func (p *lte) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	parameters = []interface{}{p.Value}
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	SQL = fmt.Sprintf("%s <= ? ", dbCol)
	return
}

//FieldName untuk keseragaman dengan interface Q. nama field
func (p *lte) FieldName() (fieldOrOperatorName string) {
	return p.Field
}

//Lte generator wrapper <=
func Lte(fieldName string, value interface{}) Q {
	return &lte{Field: fieldName, Value: value}
}
