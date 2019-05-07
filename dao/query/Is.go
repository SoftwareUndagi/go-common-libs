package query

import "fmt"

//Is untuk operator is null
type is struct {
	//Field nama struct field atau nama column database untuk query
	Field string
}

//GenerateSQL generate sql untuk Is
func (p *is) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	SQL = fmt.Sprintf("%s is null ", dbCol)
	return
}

//FieldName untuk keseragaman dengan interface Q. nama field
func (p *is) FieldName() (fieldOrOperatorName string) {
	return p.Field
}

//Is generate query wrapper for is( null)
func Is(fieldName string) Q {
	return &is{Field: fieldName}
}
