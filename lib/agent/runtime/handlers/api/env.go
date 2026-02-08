package api

import (
	"net/http"

	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/labstack/echo/v4"
)

func envList(c echo.Context) error {
	envs, err := common.ListEnvirons()
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, envs)
}
