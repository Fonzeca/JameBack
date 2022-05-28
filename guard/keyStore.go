package guard_userhub

import (
	"context"
	"errors"
	"fmt"

	"github.com/qiniu/qmgo"
	"github.com/tomogoma/go-api-guard"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDbKeyStore struct {
	db *qmgo.Database
}

func NewKeyStore(db *qmgo.Database) *MongoDbKeyStore {
	return &MongoDbKeyStore{
		db: db,
	}
}

func (ks MongoDbKeyStore) IsNotFoundError(e error) bool {
	fmt.Printf("e: %v\n", e)
	return true
}
func (ks MongoDbKeyStore) InsertAPIKey(userID string, key []byte) (api.Key, error) {
	coll := ks.db.Collection("client")
	oKey := NewApiKey(key, userID)

	c, _ := coll.Find(context.Background(), bson.M{"client": userID}).Count()
	if c > 0 {
		return nil, errors.New("Este cliente ya existe")
	}

	_, err := coll.InsertOne(context.Background(), &oKey)
	if err != nil {
		return nil, err
	}

	return oKey, nil
}
func (ks MongoDbKeyStore) APIKeyByUserIDVal(userID string, key []byte) (api.Key, error) {
	coll := ks.db.Collection("client")
	oKey := ApiKeyUserHub{}

	err := coll.Find(context.Background(), bson.M{"client": userID, "value": key}).One(&oKey)
	if err != nil {
		return nil, err
	}

	return oKey, nil
}

func (ks MongoDbKeyStore) ClientLs() ([]ApiKeyUserHub, error) {
	coll := ks.db.Collection("client")

	list := []ApiKeyUserHub{}

	err := coll.Find(context.Background(), bson.M{}).All(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
