package util

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetIntDefault(c echo.Context, param string, def int) int {
	i, e := strconv.Atoi(c.QueryParam(param))

	if e != nil {
		return def
	}

	return i
}

func GetIntRange(c echo.Context, param string, min int, max int, def int) int {
	actual := GetIntDefault(c, param, def)

	if actual > max {
		return max
	} else if actual < min {
		return min
	}

	return actual
}

func OneOf(c echo.Context, param string, options []string, def string) string {
	actual := c.QueryParam(param)

	for _, v := range options {
		if actual == v {
			return actual
		}
	}

	return def
}
