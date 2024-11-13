package data_source_test

import (
	"context"
	"testing"

	commonDataSource "bitbucket.org/moladinTech/go-lib-common/data_source"
	commonErrors "bitbucket.org/moladinTech/go-lib-common/errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionRunner(t *testing.T) {
	t.Parallel()
	dbmock, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)

	mockSqlx := sqlx.NewDb(dbmock, "sqlmock")
	tx := commonDataSource.NewTransactionRunner(mockSqlx)
	assert.Equal(t, mockSqlx, tx.DB)
}

func Test_WithTx(t *testing.T) {
	t.Parallel()
	dbmock, queryMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)

	mockSqlx := sqlx.NewDb(dbmock, "sqlmock")
	tx := commonDataSource.NewTransactionRunner(mockSqlx)
	assert.Equal(t, mockSqlx, tx.DB)

	queryMock.ExpectBegin()
	queryMock.ExpectCommit()

	successTxFunc := func(tx *sqlx.Tx) error {
		return nil
	}
	err = tx.WithTx(context.Background(), successTxFunc, nil)
	assert.Nil(t, err)
}

func Test_FailedWithTx(t *testing.T) {
	t.Parallel()
	dbmock, queryMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)

	mockSqlx := sqlx.NewDb(dbmock, "sqlmock")
	tx := commonDataSource.NewTransactionRunner(mockSqlx)
	assert.Equal(t, mockSqlx, tx.DB)

	queryMock.ExpectBegin()
	queryMock.ExpectRollback()

	successTxFunc := func(tx *sqlx.Tx) error {
		return commonErrors.ErrSQLExec
	}
	err = tx.WithTx(context.Background(), successTxFunc, nil)
	assert.ErrorContains(t, err, commonErrors.ErrSQLExec.Error())
}

func Test_StartTx(t *testing.T) {
	t.Parallel()
	dbmock, queryMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)

	queryMock.ExpectBegin()
	mockSqlx := sqlx.NewDb(dbmock, "sqlmock")
	tx, err := commonDataSource.StartTx(context.Background(), mockSqlx, nil)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
}
