package data_source

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func IsErrDuplicateKey(err error) bool {
	var psqlErr *pq.Error // ERROR postgres package has not type or method
	if errors.As(err, &psqlErr) && psqlErr.Code == "1062" {
		return true
	}

	return false
}
