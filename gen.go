package querygenerator

import (
	"fmt"
	"strings"
)

/*
**   Remove all maps from function params as they don't guarentee order of the keys
**
**
**
**
**
**
**
**
**
**
**
**
**
**
**
**/
const (
	DML_SELECT     = "SELECT"
	DML_WILDCARD   = "*"
	DML_FROM       = "FROM"
	DML_WHERE      = "WHERE"
	DML_INNER_JOIN = "INNER JOIN"
	DML_AND        = "AND"
	DML_INSERT     = "INSERT INTO"
	DML_VALUES     = "VALUES"
	DML_ON         = "ON"
)

type QueryGen struct {
	tableName string
	query     string
	values    []interface{}
	err       error
}

func GetInstance(tableName string) *QueryGen {
	return &QueryGen{tableName: tableName}
}

func (q *QueryGen) GetSelectALL() *QueryGen {

	if q.err != nil {
		return q
	}

	q.query = fmt.Sprintf("%s %s %s %s", DML_SELECT, DML_WILDCARD, DML_FROM, q.tableName)
	return q
}

func (q *QueryGen) GetSelect(columns []string) *QueryGen {

	if q.err != nil {
		return q
	}

	if len(columns) == 0 {
		q.err = fmt.Errorf("columns cannot be empty")
		return q
	}

	q.query = fmt.Sprintf("%s %s %s %s", DML_SELECT, strings.Join(columns, ","), DML_FROM, q.tableName)

	return q
}

func (q *QueryGen) WithPredicates(predicates map[string]interface{}) *QueryGen {

	if q.err != nil {
		return q
	}

	if len(predicates) == 0 {
		q.err = fmt.Errorf("predicates cannot be empty")
		return q
	}

	var p []string
	var values []interface{}

	for i, v := range predicates {
		p = append(p, fmt.Sprintf("%s=?", i))
		values = append(q.values, v)
	}

	q.query = fmt.Sprintf("%s WHERE %s", q.query, strings.Join(p, " AND "))
	q.values = values

	return q
}

func (q *QueryGen) GetInsert(args map[string]interface{}) *QueryGen {

	if q.err != nil {
		return q
	}

	if len(args) == 0 {
		q.err = fmt.Errorf("args cannot be nil")
		return q
	}

	var columns []string
	var placeholders []string
	var values []interface{}

	for i, v := range args {
		columns = append(columns, i)
		placeholders = append(placeholders, "?")
		values = append(values, v)
	}

	q.query = fmt.Sprintf("%s %s (%s) %s (%s)", DML_INSERT, q.tableName, strings.Join(columns, ","), DML_VALUES, strings.Join(placeholders, ","))
	q.values = values

	return q
}

func (q *QueryGen) WithJoins(tableName string, predicates map[string]string) *QueryGen {

	if q.err != nil {
		return q
	}

	if len(strings.TrimSpace(tableName)) <= 0 {
		q.err = fmt.Errorf("tableName cannot be empty")
		return q
	}

	if len(predicates) == 0 {
		q.err = fmt.Errorf("predicates cannot be empty")
		return q
	}

	joinQuery := fmt.Sprintf("%s %s %s", DML_INNER_JOIN, tableName, DML_ON)

	p := []string{}

	for i, v := range predicates {
		p = append(p, fmt.Sprintf("%s=%s", i, v))
	}

	if len(p) == 1 {
		joinQuery = fmt.Sprintf("%s %s", joinQuery, strings.Join(p, ""))
	} else {
		joinQuery = fmt.Sprintf(fmt.Sprintf("%s %s", joinQuery, strings.Join(p, " %s ")), DML_AND)
	}

	q.query = joinQuery

	return q
}

func (q *QueryGen) ToSQL() string {

	if q.err != nil {
		return ""
	}

	return q.query
}

func (q *QueryGen) ToSQLWithValues() (string, []interface{}) {

	if q.err != nil {
		return "", nil
	}

	return q.query, q.values
}
