package mongodb

import (
	"context"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoUserRepository struct {
	col *qmgo.Collection
}

func NewMongoUserRepository(db *qmgo.Database) domain.UserRepository {
	return &MongoUserRepository{col: db.Collection("users")}
}

func (r *MongoUserRepository) GetAll(ctx context.Context) (res []domain.User, resErr error) {
	r.col.Find(ctx, bson.M{}).All(&res)
	return res, nil
}

func (r *MongoUserRepository) Insert(ctx context.Context, user *domain.User) (res domain.User, resErr error) {
	result, err := r.col.InsertOne(ctx, user)
	if err != nil {
		return *user, err
	}

	user.Id = result.InsertedID

	return *user, nil
}
