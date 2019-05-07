package query

import "fmt"

type like struct {
	//Field nama struct field atau nama column database untuk query
	Field string
	//Value value untuk query. ini di sesuaikan dengan definisi datatype sebaiknya
	Value interface{}
	//UseNotLike flag not like
	UseNotLike bool
}

//GenerateSQL generate SQL statement untuk query >
func (p *like) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	parameters = []interface{}{p.Value}
	var dbCol string
	dbCol, err = colunmNameProvider(modelName, p.Field)
	if err != nil {
		return
	}
	if p.UseNotLike {
		SQL = fmt.Sprintf("%s not like ? ", dbCol)
	} else {
		SQL = fmt.Sprintf("%s like ? ", dbCol)
	}

	return
}

//Like generate like query operator
func Like(fieldName string, fieldValue interface{}) (query Q) {
	return &like{Field: fieldName, Value: fieldValue}
}
