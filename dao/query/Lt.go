package query

import "fmt"

//lt query dengan operator <
type lt struct {
	//Field nama struct field atau nama column database untuk query
	Field string
	//Value value untuk query. ini di sesuaikan dengan definisi datatype sebaiknya
	Value interface{}
}

//FieldName untuk keseragaman dengan interface Q. nama field
func (p *lt) FieldName() (fieldOrOperatorName string) {
	return p.Field
}

//GenerateSQL generate SQL statement untuk query <
func (p *lt) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	parameters = []interface{}{p.Value}
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	SQL = fmt.Sprintf("%s < ? ", dbCol)
	return
}

//Lt generate wrapper query for <
func Lt(fieldName string, value interface{}) Q {
	return &lt{Field: fieldName, Value: value}
}
