package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ddb "github.com/projects/cmyk-tools/handlers/db"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

var db ddb.UsersTableRepo

func init() {
	ctx := context.TODO()
	region := os.Getenv("AWS_REGION")
	logger := util.NewDevLogger(zerolog.InfoLevel)
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("handler", "cognito-auth-challenge").
			Str("gitsha", "sha...")

	})

	logger.WithContext(ctx) // add the logger to ctx so we can retrieve it
	db = ddb.NewUsersTableRepo(ctx, region)
}

type CognitoPostSignUpFn func(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error)
type cognitoPostSignUpHandler struct {
	clock      util.Clock
	usersTable string
}

func (h *cognitoPostSignUpHandler) Handler(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {

	if strings.EqualFold(event.TriggerSource, "PostConfirmation_Confirm_SignUp") {
		err := db.AddUser(context.TODO(),
			event.Request.UserAttributes["username"],
			event.Request.UserAttributes["email"],
		)

		return event, err
	}

	return event, nil
}

func NewCognitoPostSignUpHandler(clock util.Clock) CognitoPostSignUpFn {
	h := &cognitoPostSignUpHandler{
		clock: clock,
	}

	return h.Handler
}
func main() {
	lambda.Start(NewCognitoPostSignUpHandler(util.NewRealClock()))
}
