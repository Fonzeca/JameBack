package usecase

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/Fonzeka/Jame/src/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
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

func parseToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	secret := viper.GetString("jwt.secret")

	return []byte(secret), nil
}

func CheckLogged(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.Contains(c.Path(), "login") && !strings.Contains(c.Path(), "public") {

			authorization := c.Request().Header.Get("Authorization")

			re := regexp.MustCompile("Bearer (.+)")

			if !re.MatchString(authorization) {
				return utils.ErrUnauthorized
			}

			recivedToken := re.FindStringSubmatch(authorization)[1]

			token, err := jwt.Parse(recivedToken, parseToken)
			if err != nil {
				return utils.ErrExpiredToken
			}

			if !token.Valid {
				return utils.ErrExpiredToken
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Set("claims", claims)
			} else {
				return utils.ErrUnauthorized
			}
		}
		return next(c)
	}
}
