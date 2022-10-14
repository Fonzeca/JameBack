package mongodb

import (
	"context"
	"fmt"

	"github.com/Fonzeca/UserHub/server/domain"
	"github.com/Fonzeca/UserHub/server/utils"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	opt "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	col *qmgo.Collection
}

func NewMongoUserRepository(db *qmgo.Database) domain.UserRepository {
	collection := db.Collection("users")

	//Creamos el indice para que el userName no sea duplicado
	collection.CreateOneIndex(context.TODO(), options.IndexModel{Key: []string{"userName"}, IndexOptions: &opt.IndexOptions{Unique: &[]bool{true}[0]}})

	return &MongoUserRepository{col: collection}
}

func (r *MongoUserRepository) GetAll(ctx context.Context) (res []domain.User, resErr error) {
	r.col.Find(ctx, bson.M{}).All(&res)
	return res, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *domain.User) error {
	result, err := r.col.Upsert(ctx, bson.M{"userName": user.UserName}, user)
	fmt.Println(result)
	return err
}

func (r *MongoUserRepository) Delete(ctx context.Context, UserName string) error {
	return r.col.Remove(ctx, bson.M{"userName": UserName})
}

func (r *MongoUserRepository) Insert(ctx context.Context, user *domain.User) (res domain.User, resErr error) {
	result, err := r.col.InsertOne(ctx, user)
	if err != nil {
		return *user, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID)

	return *user, nil
}

func (r *MongoUserRepository) GetByUserName(ctx context.Context, userName string) (domain.User, error) {
	user := domain.User{}

	err := r.col.Find(ctx, bson.M{"userName": userName}).One(&user)

	if err != nil {
		return user, utils.ErrUserNotFound
	}

	return user, nil
}

func (r *MongoUserRepository) GetFCMTokensByUserNames(ctx context.Context, userNames []string) ([]struct {
	FCMToken string `bson:"FCMToken"`
}, error) {

	var tokens []struct {
		FCMToken string `bson:"FCMToken"`
	}

	err := r.col.Find(ctx, bson.M{"userName": bson.M{"$in": userNames}}).Select(bson.M{"FCMToken": 1}).All(&tokens)

	if err != nil {
		return tokens, utils.ErrUserNotFound
	}

	return tokens, err
}
