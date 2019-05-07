package query

import "fmt"

//InQuery query dengan operator in
type InQuery struct {
	//Field nama field
	Field string
	//Values values untuk query
	Value []interface{}
}

//GenerateSQL generate sql untuk >=
func (p *InQuery) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	parameters = []interface{}{p.Value}
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	SQL = fmt.Sprintf("%s in ? ", dbCol)
	return
}

//FieldName untuk keseragaman dengan interface Q. nama field
func (p *InQuery) FieldName() (fieldOrOperatorName string) {
	return p.Field
}

//In query dengan in
func In(field string, value []interface{}) Q {
	return &InQuery{Field: field, Value: value}
}

//InVargs dengan argument variadic
func InVargs(field string, values ...interface{}) Q {
	valueMap := make([]interface{}, len(values))
	for index, val := range values {
		valueMap[index] = val
	}
	return In(field, valueMap)
}
