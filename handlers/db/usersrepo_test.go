package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoreAndRetrieveUser(t *testing.T) {

	ctx := context.TODO()
	err := godotenv.Load(fmt.Sprintf("../../.env.local"))
	require.NoError(t, err)
	region := requiredEnvironmentVariables(t)

	repo, err := NewUsersTableRepo(ctx, region)
	require.NoError(t, err)

	u := util.RandomUser()
	savedUser, err := repo.AddTestUser(ctx, u, util.Short)
	assert.NoError(t, err, "nope")

	got, err := repo.GetUserByID(ctx, u.Id)
	assert.NoError(t, err)

	assert.Equal(t, savedUser.Id, got.Id)
	assert.Equal(t, savedUser.Email, got.Email)
	assert.Equal(t, savedUser.Name, got.Name)
	assert.Equal(t, savedUser.CreatedAt, got.CreatedAt)
	assert.Equal(t, *(savedUser.Ttl), *(got.Ttl))
}

func requiredEnvironmentVariables(t *testing.T) string {
	region := util.GetOSEnvOrFail(t, "AWS_REGION")
	_ = util.GetOSEnvOrFail(t, "USERS_TABLE")
	return region
}
