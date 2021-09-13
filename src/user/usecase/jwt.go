package usecase

import (
	"time"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type jwtCustomClaims struct {
	UserName string `json:"userName"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

func generateToken(user *domain.User) (string, error) {
	//Buscamos en los roles el index de "admin"
	isAdmin := false
	for _, v := range user.Roles {
		if v == "admin" {
			isAdmin = true
		}
	}

	exiprationMinutes := viper.GetDuration("jwt.expiration")

	//Armamos los claims
	claims := &jwtCustomClaims{
		UserName: user.UserName,
		Admin:    isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * exiprationMinutes).Unix(),
		},
	}

	// Creamos el token con las claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := viper.GetString("jwt.secret")

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
