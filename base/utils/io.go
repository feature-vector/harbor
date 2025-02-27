package utils

import (
	"io"
)

func CloseSilent(cls io.Closer) {
	_ = cls.Close()
}
