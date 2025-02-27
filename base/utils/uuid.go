package utils

import (
	"github.com/gofrs/uuid"
)

func UUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return u.String()
}
