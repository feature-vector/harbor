package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func ReadIntQuery(c *gin.Context, key string, defaults ...int) int {
	if len(defaults) > 1 {
		panic("defaults len > 1")
	}
	str := c.Query(key)
	val64, err := strconv.ParseInt(str, 10, 0)
	val := int(val64)
	if err != nil && len(defaults) > 0 {
		val = defaults[0]
	}
	return val
}

func ReadInt64Query(c *gin.Context, key string, defaults ...int64) int64 {
	if len(defaults) > 1 {
		panic("defaults len > 1")
	}
	str := c.Query(key)
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil && len(defaults) > 0 {
		val = defaults[0]
	}
	return val
}

func ReadStringQuery(c *gin.Context, key string, defaults ...string) string {
	if len(defaults) > 1 {
		panic("defaults len > 1")
	}
	str := c.Query(key)
	if str == "" && len(defaults) > 0 {
		str = defaults[0]
	}
	return str
}

func ReadTimeQuery(c *gin.Context, key string, defaults ...time.Time) time.Time {
	if len(defaults) > 1 {
		panic("defaults len > 1")
	}
	str := c.Query(key)
	if str == "" && len(defaults) > 0 {
		return defaults[0]
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return defaults[0]
	}
	return t
}

// param

func ReadIntParam(c *gin.Context, key string, defaults ...int) int {
	if len(defaults) > 1 {
		panic("defaults len > 1")
	}
	str := c.Param(key)
	val64, err := strconv.ParseInt(str, 10, 0)
	val := int(val64)
	if err != nil && len(defaults) > 0 {
		val = defaults[0]
	}
	return val
}

func ReadInt64Param(c *gin.Context, key string, defaults ...int64) int64 {
	if len(defaults) > 1 {
		panic("defaults len > 1")
	}
	str := c.Param(key)
	val, err := strconv.ParseInt(str, 10, 0)
	if err != nil && len(defaults) > 0 {
		val = defaults[0]
	}
	return val
}

func ReadStringParam(c *gin.Context, key string, defaults ...string) string {
	if len(defaults) > 1 {
		panic("defaults len > 1")
	}
	str := c.Param(key)
	if str == "" && len(defaults) > 0 {
		str = defaults[0]
	}
	return str
}
