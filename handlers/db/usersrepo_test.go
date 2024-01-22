package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/projects/cmyk-tools/handlers/model"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStoreAndRetrieveUser(t *testing.T) {

	ctx := context.TODO()
	err := godotenv.Load(fmt.Sprintf("../../.env.local"))
	require.NoError(t, err)
	region := requiredEnvironmentVariables(t)

	repo, err := NewUsersTableRepo(ctx, region)
	require.NoError(t, err)

	now := time.Now()
	u := util.RandomTestUser(util.WithCreatedAt(now))
	savedUser, err := repo.AddTestUser(ctx, u, model.Short)
	assert.NoError(t, err, "nope")

	got, err := repo.GetUserByID(ctx, u.Id)
	assert.NoError(t, err)

	assert.Equal(t, savedUser.Id, got.Id)
	assert.Equal(t, savedUser.Email, got.Email)
	assert.Equal(t, savedUser.Name, got.Name)
	assert.Equal(t, savedUser.CreatedAt, got.CreatedAt)
	assert.Equal(t, *(savedUser.MetaData.ExpiresAt), *(got.MetaData.ExpiresAt))
}

func requiredEnvironmentVariables(t *testing.T) string {
	region := util.GetOSEnvOrFail(t, "AWS_REGION")
	_ = util.GetOSEnvOrFail(t, "USERS_TABLE")
	return region
}
