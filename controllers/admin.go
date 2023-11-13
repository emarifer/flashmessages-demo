package controllers

import (
	"html/template"
	"net/http"

	"github.com/emarifer/flashmessages-demo/messages"
	"github.com/labstack/echo/v4"
)

func Admin() echo.HandlerFunc {
	return func(c echo.Context) error {
		tmpl, err := template.ParseGlob("templates/*.html")
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				err.Error(),
			)
		}

		// obtenemos el usuario logeado desde la cookie
		userCookie, _ := c.Cookie("user")

		data := map[string]any{
			"Title":    "Admin Page | Flash Messages Demo",
			"User":     userCookie.Value,
			"messages": messages.Get(c, "message"),
		}

		if err := tmpl.ExecuteTemplate(c.Response().Writer, "admin.html", data); err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				err.Error(),
			)
		}

		return nil

		/* return c.String(
			http.StatusOK,
			fmt.Sprintf(
				"Hi %s! You have been authenticated!", userCookie.Value,
			),
		) */
	}
}
