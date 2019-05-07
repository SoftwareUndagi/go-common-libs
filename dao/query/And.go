package query

import "strings"

//And query for and
type and []Q

//GenerateSQL generate sql statement
//SQL = query yang di generate. bisa berisi ?
//parameters = paramer untuk ?. urutan harus sama dengan
func (p and) GenerateSQL(modelName string, colunmNameProvider GetActualColumnName) (SQL string, parameters []interface{}, err error) {
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

	SQL = strings.Join(rsltContainer, " and ")
	return
}

//And nd query
func And(queries ...Q) Q {
	orRtvl := make(or, len(queries))
	for index, q := range queries {
		orRtvl[index] = q
	}
	return &orRtvl
}
