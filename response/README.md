# Response

## Introduction
This package is used for response format
What's got in this package.
1. Response - struct response format
2. Variable Status - [Success, Fail, Error]

## Using Package

### Using Response
```go
    // With Variable Status Success
	c.JSON(http.StatusOK, response.Response{
		Status:  StatusSuccess,
		Data:    myData,
		Message: "ok",
	})

    // With Variable Status Error
	c.JSON(http.StatusInternalServerError, response.Response{
		Status:  StatusError,
		Message: "Internal Server Error",
	})
```

### Using HttpResp

```go
    // userrepo.go
    type userRepo struct {
        commonRegistry  common.IRegistry
        masterDb *sqlx.DB
    }

    func (h *userRepo) Delete(c *gin.Context) {
        var (
            ErrSQLQueryBuilder = errors.New("error query builder") // use your own error builder package
            ctx = context.Background()
        )

        err := e.masterDb.ExecContext(ctx, query, params...)
        if err != nil {
            // will wrap error from squirrel.ToSql() with ErrSQLQueryBuilder
            return nil, 0, liberrors.WrapWithErr(err, ErrSQLQueryBuilder) // will return error to client and send notify to slack if error >= 500 
            // OR
            return nil, 0, liberrors.WrapWithErr(err, ErrSQLQueryBuilder).WithNotify(ctx, commonRegistry) // will return error to client and force to send notify
            // OR
            return nil, 0, liberrors.WrapWithErr(err, ErrSQLQueryBuilder).WithSuccessResp() // force return error to client
        }
    }

```

```go
   // userhttp.go

   response.HttpResp(ctx, err, response.ParamHttpErrResp{GinCtx: c, Registry: h.common}).
        Return(http.StatusOK, response.Response{
            Status:  response.StatusSuccess,
            Message: http.StatusText(http.StatusOK),
            Data:    map[string]interface{}{"data": "ok"},
        })
```