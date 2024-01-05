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
var logger = util.NewDevLogger(zerolog.InfoLevel)

func init() {
	db = ddb.NewUsersTableRepo(context.TODO(), os.Getenv("AWS_REGION"))
}

type CognitoPostSignUpFn func(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error)
type cognitoPostSignUpHandler struct {
	clock      util.Clock
	usersTable string
}

func (h *cognitoPostSignUpHandler) Handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {

	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("handler", "confirm-user-signup")
	})

	logger.WithContext(ctx) // add the logger to ctx so we can retrieve it
	logger.Info().Msg("handling PostConfirmation_Confirm_SignUp event")

	if strings.EqualFold(event.TriggerSource, "PostConfirmation_Confirm_SignUp") {
		_, err := db.AddUser(ctx,
			event.Request.UserAttributes["username"],
			event.Request.UserAttributes["email"],
		)

		if err != nil {
			logger.Err(err).Msg("error processing PostConfirmation_Confirm_SignUp")
		}

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
