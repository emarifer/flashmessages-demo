package controllers

import (
	"html/template"
	"net/http"

	"github.com/emarifer/flashmessages-demo/auth"
	"github.com/emarifer/flashmessages-demo/messages"
	"github.com/emarifer/flashmessages-demo/user"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignInForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		tmpl, err := template.ParseGlob("templates/*.html")
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				err.Error(),
			)
		}

		data := map[string]any{
			"Title":  "Login | Flash Messages Demo",
			"errors": messages.Get(c, "error"),
		}

		if err := tmpl.ExecuteTemplate(c.Response().Writer, "sign-in.html", data); err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				err.Error(),
			)
		}

		return nil
	}
}

func SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		// cargamos el usuario "Test User" de la DB ficticia
		storedUser := user.LoadTestUser()

		// iniciamos un nuevo struct User
		u := new(user.User)

		// parseamos el formulario submiteado de la ruta "/user/signin"
		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				err.Error(),
			)
		}

		// comparamos la contraseña hash almacenada con la versión hash
		// de la contraseña que se recibió
		if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(u.Password)); err != nil {
			/* return echo.NewHTTPError(
				http.StatusUnauthorized,
				"Password is incorrect",
			) */
			// Si las dos contraseñas no coinciden, establece
			// un mensaje y recarga la página.
			messages.Set(c, "error", "Password is incorrect!")
			c.Redirect(http.StatusMovedPermanently, c.Echo().
				Reverse("userSignInForm"))
		}

		// si la contraseña es correcta, genera tokens y configura cookies
		err := auth.GenerateTokensAndSetCookies(storedUser, c)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"Token is incorrect",
			)
		}

		messages.Set(c, "message", "Password is correct, you have been authenticated!")

		return c.Redirect(http.StatusMovedPermanently, "/admin")
	}
}
