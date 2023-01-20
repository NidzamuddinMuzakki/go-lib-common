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
