package query

//GetActualColumnName reader nama column. di pecah agar tidak terjadi cyrcular ref
type GetActualColumnName func(modelName string, name string) (colunmName string, err error)

//Q base definition untuk query
type Q interface {

	//generateSQL generate sql statement
	//SQL = query yang di generate. bisa berisi ?
	//parameters = paramer untuk ?. urutan harus sama dengan
	GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error)
}
