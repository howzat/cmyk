package confirm_user_signup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cenkalti/backoff/v4"
	"github.com/joho/godotenv"
	ddb "github.com/projects/cmyk-api/handlers/db"
	"github.com/projects/cmyk-api/handlers/model"
	"github.com/projects/cmyk-api/handlers/util"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"text/template"
	"time"
)

func TestCognitoPostSignUp_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx := context.TODO()
	err := godotenv.Load("../../../.env")
	require.NoError(t, err)

	region, userPoolID := requiredEnvironmentVariables(t)

	repo, err := ddb.NewUsersTableRepo(ctx, region)
	require.NoError(t, err)
	require.NotNil(t, repo)

	now := time.Now()
	clock := util.NewFixedClock(now)
	var handler CognitoPostSignUpFn = NewCognitoPostSignUpHandler(clock, *repo, WithLogger(util.NewDevLogger(zerolog.TraceLevel)))
	user := util.RandomTestUser()

	_, err = handler(context.TODO(), *createCognitoPostSignUpEvent(user, region, userPoolID))
	assert.NoError(t, err)

	found, err := lookupUserInDynamoDB(user.Email)
	assert.NoError(t, err)
	assert.EqualValues(t, user.Email, found.Email)
	assert.EqualValues(t, user.Id, found.Id)
	assert.NotEmpty(t, found.Id)
	assert.EqualValues(t, clock.Now(), found.CreatedAt)
}

func requiredEnvironmentVariables(t *testing.T) (string, string) {
	region := util.GetOSEnvOrFail(t, "AWS_REGION")
	userPoolID := util.GetOSEnvOrFail(t, "COGNITO_USER_POOL_ID")
	_ = util.GetOSEnvOrFail(t, "USERS_TABLE")
	return region, userPoolID
}

func lookupUserInDynamoDB(email string) (*model.User, error) {

	ctx := context.TODO()
	log := util.NewDevLogger(zerolog.TraceLevel)
	repo, err := ddb.NewUsersTableRepo(ctx, os.Getenv("AWS_REGION"))
	if err != nil {
		return nil, err
	}

	backOff := backoff.NewExponentialBackOff()
	backOff.InitialInterval = 1 * time.Second
	backOff.MaxElapsedTime = 1 * time.Minute

	var foundUser *model.User

	retryable := func() error {
		foundUser, err = repo.GetUserByID(ctx, email)
		if len(foundUser.Id) == 0 && err == nil {
			return errors.New(fmt.Sprintf("user with email not found [%s]", email))
		}

		return err
	}

	notify := func(err error, tm time.Duration) {
		if err != nil {
			log.Err(err).Msgf("error happened at time: %v", tm)
		} else {
			log.Err(err).Msgf("retry attempt complete at time: %v", tm)
		}
	}

	err = backoff.RetryNotify(retryable, backOff, notify)
	return foundUser, err
}

func createCognitoPostSignUpEvent(user model.User, region, userpoolID string) *events.CognitoEventUserPoolsPostConfirmation {

	var rawJson bytes.Buffer
	err := Create("jsonEvent", `{
        "version": "1",
        "region": "{{.Region}}",
        "userPoolId": "{{.UserpoolID}}",
        "userName": "{{.Id}}",
        "triggerSource": "PostConfirmation_ConfirmSignUp",
        "request": {
            "userAttributes": {
                "sub": "{{.Id}}",
                "cognito:email_alias": "{{.Email}}",
                "cognito:user_status": "CONFIRMED",
                "email_verified": "false",
                "email": "{{.Email}}",
                "name": "{{.Name}}"
            }
        },
        "response": {}
    }`).Execute(&rawJson, map[string]string{
		"Region":     region,
		"UserpoolID": userpoolID,
		"Id":         user.Id,
		"Email":      user.Email,
		"Name":       user.Name,
	})

	logger := util.NewDevLogger(zerolog.InfoLevel)
	logger.Info().Msg(rawJson.String())
	if err != nil {
		panic(err)
	}

	var event events.CognitoEventUserPoolsPostConfirmation
	if err := json.Unmarshal(rawJson.Bytes(), &event); err != nil {
		panic(err.Error())
	}

	return &event
}

var Create = func(name, t string) *template.Template {
	return template.Must(template.New(name).Parse(t))
}
