package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	ddb "github.com/projects/cmyk-tools/handlers/db"
	"github.com/projects/cmyk-tools/handlers/lambda/confirm-user-signup"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/rs/zerolog"
	"os"
)

var usersRepo ddb.UsersRepo

func init() {
	repo, err := ddb.NewUsersTableRepo(context.TODO(), os.Getenv("AWS_REGION"))
	if err != nil {
		panic(err)
	}
	usersRepo = *repo
}

func main() {
	lambda.Start(confirm_user_signup.NewCognitoPostSignUpHandler(
		util.NewRealClock(),
		usersRepo,
		confirm_user_signup.WithLogger(util.NewProdLogger(zerolog.InfoLevel)),
	))
}
