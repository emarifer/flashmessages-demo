package main

import (
	"net/http"

	"github.com/emarifer/flashmessages-demo/auth"
	"github.com/emarifer/flashmessages-demo/controllers"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Static("/", "assets")

	// e.Use(middleware.Logger())

	// Redireccionamos la ruta raíz a la ruta "/admin"
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/admin")
	})

	e.GET("/user/signin", controllers.SignInForm()).Name = "userSignInForm"
	e.POST("/user/signin", controllers.SignIn())

	// Defining of the admin router group.
	adminGroup := e.Group("/admin")
	// Añadimos al grupo un middlware built-in (que viene con Echo)
	// para proteger este grupo de rutas con el JWT. VER:
	// https://echo.labstack.com/docs/middleware/jwt
	adminGroup.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.Claims)
		},
		SigningKey:   []byte(auth.GetJWTSecret()),
		TokenLookup:  "cookie:access-token",
		ErrorHandler: auth.JWTErrorChecker,
	}))

	// también adjuntamos el middleware de actualización de tokens
	adminGroup.Use(auth.TokenRefresherMiddleware)

	// Router for "/admin" path.
	adminGroup.GET("", controllers.Admin())

	// Starting the server.
	e.Logger.Fatal(e.Start(":8000"))
}

/* DOCS DE ECHO:
https://echo.labstack.com/

REFERENCIAS DEL PROYECTO. VER:
https://webdevstation.com/posts/user-authentication-with-go-using-jwt-token/
https://webdevstation.com/posts/how-to-show-flash-messages-in-go-echo/

RESPECTIVOS REPOSITORIOS. VER:
https://github.com/alexsergivan/blog-examples/tree/master/authentication
https://github.com/alexsergivan/blog-examples/tree/master/flashmessages

REFERENCIAS SOBRE JWT EN ECHO, GOLANG-JWT(JWT-GO) &
SERVE STATIC FILES FROM THE PROVIDED ROOT DIRECTORY . VER:
https://echo.labstack.com/docs/middleware/jwt
https://echo.labstack.com/docs/cookbook/jwt
https://pkg.go.dev/github.com/golang-jwt/jwt/v5
https://echo.labstack.com/docs/static-files#using-echostatic

OTRAS REFERENCIAS. VER:
Mensajes Flash simples en Go:
https://www.alexedwards.net/blog/simple-flash-messages-in-golang

TAG "form" EN STRUCTS. VER:
https://stackoverflow.com/questions/68552641/is-form-an-acceptable-struct-tag-to-use-when-parsing-url-query-parameters-into

FLASH MESSAGES EN GENERAL. VER:
https://www.google.com/search?q=golang+flash+messages&oq=golang+flash+&aqs=chrome.1.69i57j0i19i512j69i60.14672j0j7&sourceid=chrome&ie=UTF-8#ip=1
https://www.google.com/search?q=go+eco+flash+message&oq=go+e&aqs=chrome.1.69i57j35i39j69i60l3j69i65l3.57312j0j4&sourceid=chrome&ie=UTF-8#ip=1

https://github.com/skynet0590/sessions/blob/v3.0.0/_examples/flash-messages/main.go
https://github.com/sujit-baniya/flash
*/
