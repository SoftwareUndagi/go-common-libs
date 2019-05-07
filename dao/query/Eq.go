package query

import "fmt"

//eq dengan operator =
type eq struct {
	//Field nama struct field atau nama column database untuk query
	Field string
	//Value value untuk query. ini di sesuaikan dengan definisi datatype sebaiknya
	Value interface{}
}

//FieldName untuk keseragaman dengan interface Q. nama field
func (p *eq) FieldName() (fieldOrOperatorName string) {
	return p.Field
}

//GenerateSQL generate SQL statement untuk query =
func (p *eq) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	parameters = []interface{}{p.Value}
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	SQL = fmt.Sprintf("%s = ? ", dbCol)
	return
}

//Eq eq builder. generate struct Eq with return interface query
func Eq(fieldName string, fieldValue interface{}) (query Q) {
	return &eq{Field: fieldName, Value: fieldValue}
}
