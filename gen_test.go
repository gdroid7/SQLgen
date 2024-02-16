package querygenerator

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

//GET INSTANCE

func TestQueryGen_GetInstance(t *testing.T) {
	temp := &QueryGen{}
	q := GetInstance("tableName")

	if reflect.TypeOf(q) != reflect.TypeOf(temp) {
		t.Errorf("Instance type does not match Expected : %s, Got: %s", reflect.TypeOf(temp), reflect.TypeOf(q))
	}

	if q.tableName != "tableName" {
		t.Errorf("Table name does not match. Expected: %s, Got: %s", "tableName", q.tableName)
	}
}

//GET SELECT ALL

func TestQueryGen_GetSelectALL_ErrorNotNIL(t *testing.T) {
	queryGen := GetInstance("tableName")
	testError := fmt.Errorf("Test error")
	queryGen.err = testError
	queryGen = queryGen.GetSelectALL()

	if queryGen.err == nil || strings.Compare(testError.Error(), queryGen.err.Error()) != 0 {
		t.Errorf("Error doesn't match. Expected: %s, Got: %s", testError.Error(), queryGen.err.Error())
	}

}

func TestQueryGen_GetSelectALL(t *testing.T) {
	// Create an instance of QueryGen
	queryGen := GetInstance("tableName")

	query := queryGen.GetSelectALL().ToSQL()

	// Assert the expected query and values
	expectedQuery := "SELECT * FROM tableName"

	if query != expectedQuery {
		t.Errorf("Query does not match. Expected: %s, Got: %s", expectedQuery, query)
	}

}

// GET SELECT
func TestQueryGen_GetSelect_WithEmptyColumns(t *testing.T) {

	queryGen := GetInstance("tableName")

	columns := []string{"column1,column2"}

	queryGen = queryGen.GetSelect(columns)

	expectedQuery := "SELECT column1,column2 FROM tableName"

	if queryGen.query != expectedQuery {
		t.Errorf("Query does not match. Expected: %s, Got: %s", expectedQuery, queryGen.query)
	}
	// Create an instance of QueryGen

	// Call the GetSelect function with test data
	columns = []string{}
	queryGen.query = ""
	// predicates := map[string]interface{}{"foo": "bar"}
	queryGen = queryGen.GetSelect(columns)

	// Assert the expected query and values
	expectedQuery = ""
	expectedError := fmt.Errorf("columns cannot be empty")

	if queryGen.query != expectedQuery {
		t.Errorf("Query does not match. Expected: %s, Got: %s", expectedQuery, queryGen.query)
	}

	if queryGen.err.Error() != expectedError.Error() {
		t.Errorf("Error does not match. Expected: %s, Got: %s", expectedError.Error(), queryGen.err.Error())
	}

}
func TestQueryGen_GetSelect_WithEmptyColumns_ErrorNotNIL(t *testing.T) {

	queryGen := GetInstance("tableName")
	testError := fmt.Errorf("Test error")
	queryGen.err = testError
	queryGen = queryGen.GetSelect([]string{})

	if queryGen.err == nil || strings.Compare(testError.Error(), queryGen.err.Error()) != 0 {
		t.Errorf("Error doesn't match. Expected: %s, Got: %s", testError.Error(), queryGen.err.Error())
	}

}

func TestQueryGen_GetSelect_WithColumns(t *testing.T) {
	// Create an instance of QueryGen
	queryGen := GetInstance("tableName")

	columns := []string{"column1,column2"}

	queryGen = queryGen.GetSelect(columns)

	expectedQuery := "SELECT column1,column2 FROM tableName"

	if queryGen.query != expectedQuery {
		t.Errorf("Query does not match. Expected: %s, Got: %s", expectedQuery, queryGen.query)
	}

}

func TestQueryGen_GetSelect_WithPredicates(t *testing.T) {
	// Create an instance of QueryGen
	queryGen := GetInstance("tableName")

	// Call the GetSelect function with test data
	columns := []string{"column1", "column2"}
	predicates := map[string]interface{}{}

	expectedError := fmt.Errorf("predicates cannot be empty")
	queryGen = queryGen.GetSelect(columns).WithPredicates(predicates)

	if queryGen.err != nil && queryGen.err.Error() != expectedError.Error() {
		t.Errorf("Error does not match. Expected: %s, Got: %s", expectedError.Error(), queryGen.err.Error())
	}

	queryGen = &QueryGen{tableName: "tableName"}
	predicates = map[string]interface{}{"foo": "bar"}

	query, values := queryGen.GetSelect(columns).WithPredicates(predicates).ToSQLWithValues()

	// Assert the expected query and values
	expectedQuery := "SELECT column1,column2 FROM tableName WHERE foo=?"
	expectedValues := []interface{}{"bar"}

	if query != expectedQuery {
		t.Errorf("Query does not match. Expected: %s, Got: %s", expectedQuery, query)
	}

	if !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("Values do not match. Expected: %v, Got: %v", expectedValues, values)
	}

}

func TestQueryGen_GetSelect_WithPredicates_ErrorNotNIL(t *testing.T) {
	queryGen := GetInstance("tableName")
	columns := []string{"column1", "column2"}
	predicates := map[string]interface{}{"foo": "bar"}

	queryGen = queryGen.GetSelect(columns)

	testError := fmt.Errorf("Test error")
	queryGen.err = testError
	queryGen = queryGen.WithPredicates(predicates)

	if queryGen.err == nil || strings.Compare(testError.Error(), queryGen.err.Error()) != 0 {
		t.Errorf("Error doesn't match. Expected: %s, Got: %s", testError.Error(), queryGen.err.Error())
	}
}

func TestQueryGen_GetInsert(t *testing.T) {

	queryGen := GetInstance("tableName")

	args := map[string]interface{}{}

	queryGen = queryGen.GetInsert(args)

	expectedError := fmt.Errorf("args cannot be nil")

	if queryGen.err == nil || strings.Compare(queryGen.err.Error(), expectedError.Error()) != 0 {
		t.Errorf("Error doesn't match. Expected: %s, Got: %s", expectedError.Error(), queryGen.err.Error())
	}

	queryGen = queryGen.GetInsert(args)

	if queryGen.err == nil || strings.Compare(queryGen.err.Error(), expectedError.Error()) != 0 {
		t.Errorf("Error doesn't match. Expected: %s, Got: %s", expectedError.Error(), queryGen.err.Error())
	}

	if query := queryGen.ToSQL(); query != "" {
		t.Errorf("SQL doensn't match. Expected: %s, Got: %s", "", query)
	}

	query, values := queryGen.ToSQLWithValues()

	if query != "" {
		t.Errorf("SQL doensn't match. Expected: %s, Got: %s", "", query)
	}

	if len := len(values); len != 0 {
		t.Errorf("SQL values length don't match. Expected: %d, Got: %d", 0, len)
	}

	queryGen = GetInstance("tableName")

	args = map[string]interface{}{
		"foo": "bar",
	}

	expectedQuery := "INSERT INTO tableName (foo) VALUES (?)"

	queryGen = queryGen.GetInsert(args)

	if strings.Compare(queryGen.query, expectedQuery) != 0 {
		t.Errorf("SQL query doensn't match. Expected: %s, Got: %s", expectedQuery, queryGen.query)
	}

	wantLen := len(args)
	gotLen := len(queryGen.values)

	if len(queryGen.values) != len(args) {
		t.Errorf("SQL values length don't match. Expected: %d, Got: %d", wantLen, gotLen)
	}

	expectedValue := "bar"

	if values := queryGen.values; values[0] != expectedValue {
		t.Errorf("SQL values order don't match. Expected: %v, Got: %v", expectedValue, values[0])
	}

}

func TestQueryGen_WithJoin(t *testing.T) {

	queryGen := GetInstance("tableName")
	predicates := map[string]string{}
	expectedError := fmt.Errorf("tableName cannot be empty")
	queryGen = queryGen.WithJoins("", predicates)

	if queryGen.err == nil {
		t.Errorf("Errors don't match. Expected: %v, Got: %v", expectedError, queryGen.err.Error())
	}

	if queryGen.err.Error() != expectedError.Error() {
		t.Errorf("Errors don't match. Expected: %v, Got: %v", expectedError, queryGen.err.Error())
	}

	queryGen = queryGen.WithJoins("", predicates)

	if queryGen.err == nil {
		t.Errorf("Errors don't match. Expected error to exist, Got: nil")
	}

	queryGen.err = nil

	queryGen = queryGen.WithJoins("tableName", predicates)

	expectedError = fmt.Errorf("predicates cannot be empty")

	if queryGen.err.Error() != expectedError.Error() {
		t.Errorf("Errors don't match. Expected: %v, Got: %v", expectedError, queryGen.err.Error())
	}

	queryGen.err = nil

	predicates = map[string]string{"foo": "bar"}

	expectedQuery := "INNER JOIN tableName ON foo=bar"

	queryGen = queryGen.WithJoins("tableName", predicates)

	if strings.Compare(expectedQuery, queryGen.query) != 0 {
		t.Errorf("SQL doesn't match. Expected: %s, Got: %s", expectedQuery, queryGen.query)
	}

	predicates = map[string]string{"foo": "bar"}

	expectedQuery = "INNER JOIN tableName ON foo=bar"

	queryGen = queryGen.WithJoins("tableName", predicates)

	if strings.Compare(expectedQuery, queryGen.query) != 0 {
		t.Errorf("SQL doesn't match. Expected: %s, Got: %s", expectedQuery, queryGen.query)
	}
}
