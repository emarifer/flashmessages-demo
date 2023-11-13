package messages

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// nombre de la cookie
const sessionName = "fmessages"

func getCookieStore() *sessions.CookieStore {
	// esto debe estar un .env file
	sessionKey := "test-session-key"

	return sessions.NewCookieStore([]byte(sessionKey))
}

// Set agrega un nuevo mensaje al almacén de cookies.
func Set(c echo.Context, name, value string) {
	session, _ := getCookieStore().Get(c.Request(), sessionName)

	session.AddFlash(value, name)

	session.Save(c.Request(), c.Response())
}

// Get recibe mensajes flash del almacén de cookies.
func Get(c echo.Context, name string) []string {
	session, _ := getCookieStore().Get(c.Request(), sessionName)

	fm := session.Flashes(name)

	// si hay algunos mensajes…
	if len(fm) > 0 {
		session.Save(c.Request(), c.Response())

		// iniciamos un strings slice vacío que luego retornamos
		// con los mensajes
		var flashes []string
		for _, fl := range fm {
			// añadimos los mensajes al slice
			flashes = append(flashes, fl.(string))
		}

		return flashes
	}

	return nil
}
