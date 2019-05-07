package query

import "fmt"

//IsNot untuk operator is not null
type isNot struct {
	//Field nama struct field atau nama column database untuk query
	Field string
}

//GenerateSQL generate sql untuk Is
func (p *isNot) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	SQL = fmt.Sprintf("%s is not null ", dbCol)
	return
}

//FieldName untuk keseragaman dengan interface Q. nama field
func (p *isNot) FieldName() (fieldOrOperatorName string) {
	return p.Field
}

//IsNot generate wrapper for is not
func IsNot(fieldName string) Q {
	return &isNot{Field: fieldName}
}
