package menunda

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (m *Menunda) routes() *echo.Echo {
	r := echo.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	if m.Debug {
		r.Use(middleware.Logger())
	}
	r.Use(middleware.Recover())

	r.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Mantap")
	})

	return r

}
