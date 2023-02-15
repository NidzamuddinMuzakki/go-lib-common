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

func Test_Exec(t *testing.T) {
	type (
		args struct {
			ctx  context.Context
			stmt *sqlx.Stmt
		}
		want struct {
			err error
		}
		scenario func(args args) (result want)
		testCase struct {
			name string
			args
			want
			scenario
		}
	)

	testCases := []testCase{
		{
			name: "Normal, empty destination",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
			scenario: func(args args) (result want) {
				dbmock, queryMock := mockSqlx()
				assert.NotNil(t, dbmock)
				query := "update from table1 set id = $1 where id = $2"

				queryMock.ExpectPrepare(query).ExpectExec().
					WillReturnResult(sqlmock.NewResult(1, 1))

				commonStmt := commonDataSource.NewStatement(nil, query, 1, 2)
				result.err = commonDataSource.Exec(args.ctx, dbmock, commonStmt)
				return result
			},
		},
		{
			name: "Normal, array destination",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
			scenario: func(args args) (result want) {
				dbmock, queryMock := mockSqlx()
				assert.NotNil(t, dbmock)
				query := "select * from table1 where id = $1"

				queryMock.ExpectPrepare(query).ExpectQuery().
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name"},
						).AddRow(
							1, "John Doe",
						),
					)

				dummyDest := []struct {
					Id   int    `db:"id"`
					Name string `db:"name"`
				}{}
				commonStmt := commonDataSource.NewStatement(&dummyDest, query, 1)
				commonStmt.Debug()
				result.err = commonDataSource.Exec(args.ctx, dbmock, commonStmt)
				return result
			},
		},
		{
			name: "Normal, primitiveVariable destination",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
			scenario: func(args args) (result want) {
				dbmock, queryMock := mockSqlx()
				assert.NotNil(t, dbmock)
				query := "select id from table1 where id = $1"

				queryMock.ExpectPrepare(query).ExpectQuery().
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id"},
						).AddRow(
							1,
						),
					)

				dummyId := 0
				commonStmt := commonDataSource.NewStatement(&dummyId, query, 1)
				commonStmt.Debug()
				result.err = commonDataSource.Exec(args.ctx, dbmock, commonStmt)
				return result
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.scenario(testCase.args)

			if testCase.want.err != nil {
				assert.ErrorContains(t, result.err, testCase.want.err.Error())
			} else {
				assert.Nil(t, result.err)
			}
		})
	}
}

func Test_ExecTx(t *testing.T) {
	type (
		args struct {
			ctx  context.Context
			stmt *sqlx.Stmt
		}
		want struct {
			err error
		}
		scenario func(args args) (result want)
		testCase struct {
			name string
			args
			want
			scenario
		}
	)

	testCases := []testCase{
		{
			name: "Normal, empty destination",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
			scenario: func(args args) (result want) {
				dbmock, queryMock := mockSqlx()
				assert.NotNil(t, dbmock)
				query := "update from table1 set id = $1 where id = $2"

				queryMock.ExpectBegin()
				queryMock.ExpectPrepare(query).ExpectExec().
					WillReturnResult(sqlmock.NewResult(1, 1))
				queryMock.ExpectCommit()

				commonStmt := commonDataSource.NewStatement(nil, query, 1, 2)
				result.err = commonDataSource.ExecTx(args.ctx, dbmock, commonStmt)
				return result
			},
		},
		{
			name: "Error and Rollback case",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: commonErrors.ErrSQLExec,
			},
			scenario: func(args args) (result want) {
				dbmock, queryMock := mockSqlx()
				assert.NotNil(t, dbmock)
				query := "update from table1 set id = $1 where id = $2"

				queryMock.ExpectBegin()
				queryMock.ExpectPrepare(query).ExpectExec().
					WillReturnError(commonErrors.ErrSQLExec)
				queryMock.ExpectRollback()

				commonStmt := commonDataSource.NewStatement(nil, query, 1, 2)
				result.err = commonDataSource.ExecTx(args.ctx, dbmock, commonStmt)
				return result
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.scenario(testCase.args)

			if testCase.want.err != nil {
				assert.ErrorContains(t, result.err, testCase.want.err.Error())
			} else {
				assert.Nil(t, result.err)
			}
		})
	}
}

func mockSqlx() (mockSqlx *sqlx.DB, queryMock sqlmock.Sqlmock) {
	dbmock, queryMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, queryMock
	}

	mockSqlx = sqlx.NewDb(dbmock, "sqlmock")
	return mockSqlx, queryMock
}
