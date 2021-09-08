package mongodb_test

import (
	"context"
	"testing"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/Fonzeka/Jame/src/user/repository/mongodb"
	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAll(t *testing.T) {

	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "Login_test", Coll: "user"})

	if err != nil {
		return
	}

	defer func() {
		if err = cli.Close(ctx); err != nil {
			panic(err)
		}
	}()
	defer cli.DropDatabase(ctx)

	repo := mongodb.NewMongoUserRepository(cli.Database)

	assert.NotNil(t, repo)

	singleUser := domain.User{
		UserName: "singleUser",
		Roles:    []string{"rol1", "rol2"},
	}

	manyUsers := []domain.User{
		{UserName: "manyUsers1"},
		{UserName: "manyUsers2"},
		{UserName: "manyUsers3"},
		{UserName: "manyUsers4"},
		{UserName: "manyUsers5"},
	}

	t.Run("insert-one", func(t *testing.T) {

		usr, err := repo.Insert(ctx, &singleUser)
		assert.NoError(t, err)
		assert.NotEmpty(t, usr.Id)
		assert.NotEmpty(t, singleUser.Id)
	})

	t.Run("insert-duplcate", func(t *testing.T) {
		singleUser.Id = primitive.NilObjectID

		usr, err := repo.Insert(ctx, &singleUser)
		assert.Error(t, err)
		assert.NotEmpty(t, usr.Id) //Se usa NotEmpty porque el driver ya le asigna el id antes de insertarlo a la BD
	})

	t.Run("get-by-username", func(t *testing.T) {
		usr, err := repo.GetByUserName(ctx, singleUser.UserName)
		assert.NoError(t, err)
		assert.Equal(t, usr.Roles, singleUser.Roles)
	})

	t.Run("error-get-by-username", func(t *testing.T) {
		_, err := repo.GetByUserName(ctx, "singleUser.UserName")
		assert.Error(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		err := repo.Delete(ctx, singleUser.UserName)
		assert.NoError(t, err)
	})

	t.Run("get-all", func(t *testing.T) {
		for _, u := range manyUsers {
			_, err := repo.Insert(ctx, &u)
			assert.NoError(t, err)
		}

		users, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, len(users), len(manyUsers))

		//TODO: Validate each element
	})

}
