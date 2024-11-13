package data_source_test

import (
	"testing"

	commonDataSource "bitbucket.org/moladinTech/go-lib-common/data_source"
	"github.com/stretchr/testify/assert"
)

func Test_Destination(t *testing.T) {
	t.Parallel()
	statement := commonDataSource.NewStatement(
		nil,
		"select * from table1 where id = $1",
		1,
	)

	var expectedDestination string
	statement = statement.SetDestination(expectedDestination)
	assert.Equal(
		t,
		expectedDestination,
		statement.GetDestination(),
	)
}

func Test_Query(t *testing.T) {
	t.Parallel()
	statement := commonDataSource.NewStatement(
		nil,
		"select * from table1 where id = $1",
		1,
	)

	var expectedQuery string = "select * from table2 where id = $1"
	statement = statement.SetQuery(expectedQuery)
	assert.Equal(
		t,
		expectedQuery,
		statement.GetQuery(),
	)
}

func Test_Args(t *testing.T) {
	t.Parallel()
	statement := commonDataSource.NewStatement(
		nil,
		"select * from table1 where id = $1",
		1,
	)

	var expectedArgs = make([]interface{}, 1)
	expectedArgs[0] = 1
	statement = statement.SetArgs(expectedArgs)
	assert.Equal(
		t,
		expectedArgs,
		statement.GetArgs(),
	)
}
