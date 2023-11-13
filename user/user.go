package user

import "golang.org/x/crypto/bcrypt"

type User struct {
	Password string `json:"password" form:"password"`
	Name     string `json:"name" form:"name"`
}

// Creamos un usuario con la contraseña de "test" cifrada, simulando un
// usuario obtenido de una DB

func LoadTestUser() *User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), 8)

	return &User{
		Password: string(hashedPassword),
		Name:     "Test User",
	}
}

/* SOBRE EL TAG "form". VER:
https://stackoverflow.com/questions/68552641/is-form-an-acceptable-struct-tag-to-use-when-parsing-url-query-parameters-into
*/
