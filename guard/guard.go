package guard_userhub

import (
	"encoding/base64"
	"fmt"

	"github.com/Fonzeca/UserHub/server/utils"
	"github.com/labstack/echo/v4"
	"github.com/tomogoma/go-api-guard"
)

type Guard struct {
	core  *api.Guard
	store *MongoDbKeyStore
}

func NewGuard(generator *KeyGeneratorUserHub, store *MongoDbKeyStore) *Guard {
	g, _ := api.NewGuard(store, api.WithKeyGenerator(generator))
	return &Guard{
		core:  g,
		store: store,
	}
}

func (g *Guard) ClientLs() ([]string, error) {
	list, err := g.store.ClientLs()
	if err != nil {
		return nil, err
	}

	var listStr []string
	for _, v := range list {
		listStr = append(listStr, v.Client)
	}

	return listStr, nil
}

func (g *Guard) GenerateAndSaveApiKey(client string) (*ApiKeyUserHub, error) {
	APIKey, err := g.core.NewAPIKey(client)
	if err != nil {
		return nil, err
	}
	keyHub := APIKey.(ApiKeyUserHub)
	return &keyHub, nil
}

func (g *Guard) EchoMiddlewareApiKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		base64ApiKey := c.Request().Header.Get("ApiKey")

		apiKey, err := base64.StdEncoding.DecodeString(base64ApiKey)
		if err != nil {
			return utils.ErrUnauthorized
		}

		client, err := g.core.APIKeyValid(apiKey)
		if err != nil {
			return utils.ErrUnauthorized
		}
		fmt.Println("Autorizado client: " + client)

		return next(c)
	}
}
