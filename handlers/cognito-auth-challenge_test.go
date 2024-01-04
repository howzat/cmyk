package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	ddb "github.com/projects/cmyk-tools/handlers/db"
	"github.com/projects/cmyk-tools/handlers/model"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"text/template"
	"time"
)

func TestCognitoPostSignUp_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := godotenv.Load("../.env", "../.env.static-refs")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	region := os.Getenv("AWS_REGION")
	require.NotEmpty(t, region, "region must not be empty")
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	require.NotEmpty(t, userPoolID, "COGNITO_USER_POOL_ID must not be empty")

	now := time.Now()
	clock := util.NewFixedClock(now)
	var handler CognitoPostSignUpFn = NewCognitoPostSignUpHandler(clock)
	user := util.RandomUser()

	_, err = handler(*createCognitoPostSignUpEvent(user, region, userPoolID))
	assert.NoError(t, err)

	found, err := lookupUserInDynamoDB(user.Email)
	assert.NoError(t, err)
	assert.EqualValues(t, user.Email, found.Email)
	assert.EqualValues(t, user.Username, found.Username)
	assert.NotEmpty(t, found.Id)
	assert.EqualValues(t, clock.Now(), found.CreatedAt)
}

func lookupUserInDynamoDB(email string) (*model.User, error) {

	ctx := context.TODO()
	repo := ddb.NewUsersTableRepo(ctx, os.Getenv("AWS_REGION"))
	return repo.GetUser(ctx, email)
}

func createCognitoPostSignUpEvent(user model.User, userpoolID, region string) *events.CognitoEventUserPoolsPostConfirmation {

	var rawJson bytes.Buffer
	err := Create("jsonEvent", `{
        "version": "1",
        "region": "{{.Region}}",
        "userPoolId": "{{.UserpoolID}}",
        "userName": "{{.Username}}",
        "triggerSource": "PostConfirmation_ConfirmSignUp",
        "request": {
            "userAttributes": {
                "sub": "{{.Username}}",
                "cognito:email_alias": "{{.Email}}",
                "cognito:user_status": "CONFIRMED",
                "email_verified": "false",
                "email": "{{.Email}}"
            }
        },
        "response": {}
    }`).Execute(&rawJson, map[string]string{
		"Region":     region,
		"UserpoolID": userpoolID,
		"Username":   user.Username,
		"Email":      user.Email,
	})

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
