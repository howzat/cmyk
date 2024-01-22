package confirm_user_signup

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	ddb "github.com/projects/cmyk-tools/handlers/db"
	"github.com/projects/cmyk-tools/handlers/model"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/rs/zerolog"
	"strings"
)

type CognitoPostSignUpFn func(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error)
type cognitoPostSignUpHandler struct {
	clock     util.Clock
	logger    zerolog.Logger
	usersRepo ddb.UsersRepo
}

func (h *cognitoPostSignUpHandler) Handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {

	logger := h.logger
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("handler", "confirm-user-signup")
	})

	logger.WithContext(ctx) // add the logger to ctx so we can retrieve it
	logger.Info().Msg("handling PostConfirmation_Confirm_SignUp event")

	if strings.EqualFold(event.TriggerSource, "PostConfirmation_ConfirmSignUp") {
		_, err := h.usersRepo.AddUser(ctx, model.User{
			Id:        event.Request.UserAttributes["sub"],
			Email:     event.Request.UserAttributes["email"],
			Name:      event.Request.UserAttributes["name"],
			CreatedAt: h.clock.Now(),
		})

		if err != nil {
			logger.Err(err).Msg("error processing PostConfirmation_Confirm_SignUp")
		}

		return event, err
	} else {
		logger.Debug().Str("TriggerSource", event.TriggerSource).Msg("event is noop for handler")
	}

	return event, nil
}

type CognitoPostSignUpHandlerOption = func(handler *cognitoPostSignUpHandler) *cognitoPostSignUpHandler

func WithLogger(logger zerolog.Logger) CognitoPostSignUpHandlerOption {
	return func(h *cognitoPostSignUpHandler) *cognitoPostSignUpHandler {
		return &cognitoPostSignUpHandler{
			clock:     h.clock,
			logger:    logger,
			usersRepo: h.usersRepo,
		}
	}
}

func NewCognitoPostSignUpHandler(clock util.Clock, usersRepo ddb.UsersRepo, options ...CognitoPostSignUpHandlerOption) CognitoPostSignUpFn {
	h := &cognitoPostSignUpHandler{
		clock:     clock,
		usersRepo: usersRepo,
	}

	for _, option := range options {
		h = option(h)
	}

	return h.Handler
}
