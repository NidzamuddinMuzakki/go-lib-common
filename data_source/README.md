# Data Source

## Introduction
This package is used for configuration database
What's got in this package.
1. NewDB - used to open connection your database.
2. GetDbColumnsAndValue - used to get value data in struct model.
3. Exec - used to wrapping multiple queries or single query without transaction.
4. ExecTx - used to wrapping multiple queries or single query in a transaction.

## Using Package

### Using NewDB
```go
	masterDbCon, err = data_source.NewDB(&commonDs.Config{
		Driver:                "Your Driver",
		Host:                  "Your Host",
		Port:                  3306, // Your Port
		DBName:                "Your DBName",
		User:                  "Your User",
		Password:              "Your Password",
		SSLMode:               "Your SSLMode",
		MaxOpenConnections:    10, // Your Max Open Connections
		MaxLifeTimeConnection: 60, // Your Max Life Time Connection
		MaxIdleConnections:    10, // Your Max Idle Connections
		MaxIdleTimeConnection: 30, // Your Max Idle Time Connection
	})
	if err != nil {
		panic(err)
	}
```

### Using GetDbColumnsAndValue
```go
    type Person struct {
        Name    string `db:"name"`
        Address string `db:"address"`
    }
    // 1. without excluded
    data := GetDbColumnsAndValue(
		Person{Name: "John", Address: "Jakarta"},
	)
    fmt.Println(data) //map[address:Pontianak name:Hakaman]

    // 2. with excluded
    data = GetDbColumnsAndValue(
		Person{Name: "John", Address: "Jakarta"},
        "address"
	)
    fmt.Println(data) //map[name:Hakaman]
```

### Using Exec
```go
    //Imagine you have data Person with {id:1, name: "John", address: "Jakarta"}
    type Person struct {
        Name    string `db:"name"`
        Address string `db:"address"`
    }

    var result Person

    err = data_source.Exec(ctx, r.master, data_source.NewStatement(
        &result,
        "SELECT name, address from person where id=$1",
        [1],
    ))
```

### Using ExecTx
```go
    //Imagine you have data Person with {id:1, name: "John", address: "Jakarta"}
    type Person struct {
        Name    string `db:"name"`
        Address string `db:"address"`
    }

    var result Person

    err = data_source.ExecTx(ctx, r.master, 
        data_source.NewStatement(
            &result,
            "SELECT name, address from person where id=$1",
            [1],
        ),
        data_source.NewStatement(
            &result,
            "SELECT name, address from person where id=$1",
            [2],
        ),
    )
```

### Using WithTx
*note: for each repository must have same database
 ```go
   // main.go
   ...
   transaction := data_source.NewTransaction(db)
   repositoryRegistry := repository.NewRegistry(db)
   ...
   // repository/registry.go
   ...
   func NewRegistry(transaction) {
     return &registry{
         transaction: transaction
     }
   }
   ...
   // service/user.go
   ...
   txFunc   := TxFunc(func(tx *sqlx.Tx) error {
     user, err := repository.user.Create(ctx, tx, user)
     if err != nil {
         return err
     }
     role, err := repository.role.Create(ctx, tx, role)
     if err != nil {
        return err
     }
   })
   err := repository.transaction.WithTx(ctx, txFunc, nil)
   if err != nil {
     return err
   }
   ...
```