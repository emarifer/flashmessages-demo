package auth

import (
	"net/http"
	"time"

	"github.com/emarifer/flashmessages-demo/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"

	// En Prod esto debe colocarse en un archivo .env
	jwtSecretKey        = "some-secret-key"
	jwtRefreshSecretKey = "some-refresh-secret-key"
)

func GetJWTSecret() string {
	return jwtSecretKey
}

func GetRefreshJWTSecret() string {
	return jwtRefreshSecretKey
}

// Creamos una estructura que se codificará en un JWT.
// Agregamos jwt.RegisteredClaims como un tipo incrustado,
// para proporcionar campos como el tiempo de vencimiento, p.ej.
type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

// La función GenerateTokensAndSetCookies genera el token jwt
// y lo guarda en una cookie "HttpOnly" (VER nota abajo).

func GenerateTokensAndSetCookies(user *user.User, c echo.Context) error {

	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)

	setUserCookie(user, exp, c)

	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return nil
}

func generateAccessToken(user *user.User) (string, time.Time, error) {
	// tiempo de vencimiento del token
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

func generateRefreshToken(user *user.User) (string, time.Time, error) {
	// tiempo de vencimiento del token
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func generateToken(user *user.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {

	// creamos un JWT Claims, que incluye el nombre de usuario
	// y el tiempo de vencimiento.
	claims := &Claims{
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		/* StandardClaims: jwt.StandardClaims{
			// En JWT, el tiempo de caducidad se expresa en milisegundos Unix.
			ExpiresAt: expirationTime.Unix(),
		}, */
	}

	// generamos el token con el algoritmo HS256 utilizado
	// en la firma y en el Claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// creamos el string JWT
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

// creamos una nueva cookie, que almacenará el token JWT válido.

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	// Http-only ayuda a mitigar el riesgo de que scripts
	// del lado del cliente accedan a la cookie protegida.
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// El propósito de esta otra cookie es almacenar el nombre del usuario.

func setUserCookie(user *user.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Name
	cookie.Expires = expiration
	cookie.Path = "/"
	// Esta cookie si puede ser accedida por el cliente,
	// porque no dejamos http-only=true
	// cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// JWTErrorChecker se ejecutará cuando el usuario
// intente acceder a una ruta protegida.

func JWTErrorChecker(c echo.Context, err error) error {
	return c.Redirect(
		http.StatusMovedPermanently,
		c.Echo().Reverse("userSignInForm"),
	)
}

// Middleware TokenRefresherMiddleware, que actualiza
// los tokens JWT si el token de acceso está a punto de caducar.

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// si el usuario no está autenticado (no hay datos del token de
		// usuario en el contexto), no hagas nada.
		if c.Get("user") == nil {
			return next(c)
		}

		// en caso contrario lo obtnemos y lo casteamos, primeramente
		// como *jwt.Token y luego como nuestro custom Claims.
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*Claims)

		// Nos aseguramos de que no se emita un nuevo token hasta que haya
		// transcurrido el tiempo suficiente.
		// En este caso, sólo se emitirá un nuevo token si el token anterior
		// está dentro un intervalo de 15 minutos antes de su caducidad.
		// fmt.Println("Tiempo hasta ahora:", time.Until(claims.ExpiresAt.Time))
		if time.Until(claims.ExpiresAt.Time) < 15*time.Minute {
			// obtenemos el token de actualización de la cookie.
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
				// analizamos el token y compruebamos si es válido.
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})

				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					// si todo está bien, actualizamos los tokens.
					_ = GenerateTokensAndSetCookies(&user.User{
						Name: claims.Name,
					}, c)
				}
			}

		}

		return next(c)
	}
}

/* FROM https://www.cookiepro.com/knowledge/httponly-cookie/

"Una cookie HttpOnly es una etiqueta agregada a una cookie del navegador que evita que los scripts del lado del cliente accedan a los datos. Proporciona una puerta que impide que cualquier otra persona que no sea el servidor acceda a la cookie especializada. El uso de la etiqueta HttpOnly al generar una cookie ayuda a mitigar el riesgo de que los scripts del lado del cliente accedan a la cookie protegida, lo que hace que estas cookies sean más seguras."

DOCS DE golang-jwt/jwt (SUSTITUYE A github.com/dgrijalva/jwt-go). VER:
https://pkg.go.dev/github.com/golang-jwt/jwt/v5
*/
