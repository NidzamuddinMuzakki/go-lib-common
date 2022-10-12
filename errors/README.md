## Why you should switch to this error wrapper package

This package was inspired from [`errors`](https://github.com/PumpkinSeed/errors) Where we have multiple error that's called from one func to another func, creating a simple stack trace that easy to read, getting the root cause from that multiple error, and append key when call wrap as response with corresponding to the key

## Using Package

### Using WrapWithErr


```go

    var (
        ErrSQLQueryBuilder = errors.New("error query builder") // use your own error builder package
    )

    sql, args, err := recon.ToSql()
	if err != nil {
        // will wrap error from recon.ToSql() with  ErrSQLQueryBuilder
		return nil, 0, liberrors.WrapWithErr(err, ErrSQLQueryBuilder)
	}


```

### Using Wrap

```go

    
    func ToSQL(obj interface{}) (string,string,error){
        return "","", errors.New("error query builder")
    }

    sql, args, err := ToSql(query)
	if err != nil {
        // will wrap error from recon.ToSql() and will get the stack trace
        // will be usefull when you need to know who is the caller
		return nil, 0, liberrors.Wrap(err, ErrSQLQueryBuilder)
	}

```

### Get Root Error / Root Cause

```go

    func ToSQL(obj interface{}) (string,string,error){
        return "","", errors.New("error query builder")
    }

    sql, args, err := ToSql(query)
	if err != nil {
        // will wrap error from recon.ToSql() and will get the stack trace
        // will be usefull when you need to know who is the caller
		errWrapped = liberrors.Wrap(err, ErrSQLQueryBuilder)
	}

    rootCauseErr := liberrors.RootCause(errWrapped)

    same := rootCauseErr == err // true

```


### Getting Key Error

```go
     var (
        ErrSQLQueryBuilder = errors.New("error query builder") // use your own error builder package
    )

    sql, args, err := recon.ToSql()
	if err != nil {
        // will wrap error from recon.ToSql() with  ErrSQLQueryBuilder
		errWrapped = liberrors.WrapWithErr(err, ErrSQLQueryBuilder)
	}

    KeyError := liberrors.GetErrKey(errWrapped)

    same := KeyError == ErrSQLQueryBuilder // true
    notSame := KeyError == err // false, since its wrapped by new errror, the key is changed to error wrapped



```

### Key Error As Key on MapResponse

```go
    // let's say we have this map and error variable

    var (
        ErrSQLQueryBuilder = errors.New("error query builder")
    )

    type ResponseError struct {
        Message string
        HTTPStatusCode int
        Error error
    }

    MapResponse := map[error]ResponseError{
        ErrSQLQueryBuilder: ResponseError{
            Message: "Failed When Executing Query Builder"
            HTTPStatusCode: http.StatusInternalServerError
        }
    }


    // and then we have this function

    sql, args, err := recon.ToSql()
	if err != nil {
        // will wrap error from recon.ToSql() with  ErrSQLQueryBuilder
		errWrapped = liberrors.WrapWithErr(err, ErrSQLQueryBuilder)
	}

    // we want the key error to get our http error response
    KeyError := liberrors.GetErrKey(errWrapped)

    // remember that everytime error is wrapped with new error, new error will be the key instead


    response := MapResponse[KeyError] // will return ResponseError{
            //Message: "Failed When Executing Query Builder"
            //HTTPStatusCode: http.StatusInternalServerError
      //  }


```

### Compare Error

```go
     var (
        ErrSQLQueryBuilder = errors.New("error query builder") // use your own error builder package
    )

    sql, args, err := recon.ToSql()
	if err != nil {
        // will wrap error from recon.ToSql() with  ErrSQLQueryBuilder
		errWrapped := liberrors.WrapWithErr(err, ErrSQLQueryBuilder)
	}

    wrap := liberrrors.Wrap(errWrapped)

    assert(true, errors.Is(wrap,wrap)) // will be using Is Interface 
    assert(true, errors.Is(wrap,errWrapped))
    assert(true, errors.Is(wrap, err))

 ```

### Printing Error

 Remember That when printing error, it will also printed the stack trace. so if you only want to printed the original / root cause error, you have to called RootErr function first

 ```go
    var (
        ErrSQLQueryBuilder = errors.New("error query builder") // use your own error builder package
    )

    sql, args, err := recon.ToSql()
	if err != nil {
        // will wrap error from recon.ToSql() with  ErrSQLQueryBuilder
		errWrapped := liberrors.WrapWithErr(err, ErrSQLQueryBuilder)
	}

    logger.Error(ctx, errWrapped.Error()) // error query builder -- At /Users/Moladin/go/go-lib-common/common/readme.md: 168: root cause: error table x is not found

```
    






