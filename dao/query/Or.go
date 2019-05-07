package query

import "strings"

//MarkerOrOperator nama field or, ini meniru sequelize
const MarkerOrOperator = "$or"

//Or or query container
type or []Q

//GenerateSQL generate sql statement
//SQL = query yang di generate. bisa berisi ?
//parameters = paramer untuk ?. urutan harus sama dengan
func (p or) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
	if len(p) == 0 {
		return
	}
	rsltContainer := []string{}
	for _, q := range p {
		curSQL, paramCur, errCurr := q.GenerateSQL(modelName, colunmNameProvider)
		if errCurr != nil {
			err = errCurr
			return
		}
		rsltContainer = append(rsltContainer, "( "+curSQL+" )")
		if len(paramCur) > 0 {
			parameters = append(parameters, paramCur...)
		}
	}

	SQL = strings.Join(rsltContainer, " or ")
	return
}

//Or generate or query
func Or(queries ...Q) Q {
	orRtvl := make(or, len(queries))
	for index, q := range queries {
		orRtvl[index] = q
	}
	return &orRtvl
}
