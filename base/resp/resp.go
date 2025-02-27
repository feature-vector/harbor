package resp

import (
	"errors"
	"github.com/feature-vector/harbor/base/env"
	"github.com/feature-vector/harbor/base/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Success(c *gin.Context, data ...interface{}) {
	if len(data) > 1 {
		panic("resp.Success only support 0 or 1 param")
	}
	if len(data) == 1 {
		c.JSON(http.StatusOK, data[0])
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func Error(c *gin.Context, err error) {
	r := respFromErr(err)
	c.JSON(r.HttpCode, r)
}

func Abort(c *gin.Context, err error) {
	r := respFromErr(err)
	c.AbortWithStatusJSON(r.HttpCode, r)
}

func respFromErr(err error) errs.ProtocolError {
	pe := errs.ProtocolError{}
	if errors.As(err, &pe) {
		return pe
	}
	pe.HttpCode = http.StatusInternalServerError
	pe.Code = errs.ErrCodeUnknownError
	if env.IsProduction() {
		pe.Message = "Server is busy now"
	} else {
		pe.Message = err.Error()
	}
	zap.L().Error(
		"unknown server err",
		zap.Error(err),
	)
	return pe
}
