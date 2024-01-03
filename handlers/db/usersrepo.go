package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/projects/cmyk-tools/handlers/model"
	"github.com/projects/cmyk-tools/handlers/util"
	"log"
	"time"
)

type UsersTableRepo struct {
	ddb   DynamoRepository
	clock util.Clock
}

func NewUsersTableRepo(ctx context.Context, region string) UsersTableRepo {
	return UsersTableRepo{
		ddb:   NewInstance(ctx, region, UsersTableEnvKey),
		clock: util.NewRealClock(),
	}
}

func (r *UsersTableRepo) AddUser(ctx context.Context, username, email string) error {

	currentTime, id, err := util.CurrentTimeAndULID(r.clock)
	if err != nil {
		log.Fatal(err)
	}

	user := createUser{
		Username:  username,
		Email:     email,
		Pk:        id.String(),
		Id:        id.String(),
		CreatedAt: currentTime.Format(time.RFC3339),
	}

	idPK, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	user.Pk = user.Email
	emailPK, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.ddb.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					Item:                idPK,
					TableName:           aws.String(r.ddb.Tablename),
					ConditionExpression: aws.String("attribute_not_exists(pk)"),
				},
			},
			{
				Put: &types.Put{
					Item:                emailPK,
					TableName:           aws.String(r.ddb.Tablename),
					ConditionExpression: aws.String("attribute_not_exists(pk)"),
				},
			},
		}})

	return err
}

func (r *UsersTableRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.ddb.GetByKey(ctx, map[string]types.AttributeValue{
		"email": &types.AttributeValueMemberS{Value: email},
	}, &user)

	return &user, err
}

type createUser struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Pk        string `json:"pk" validate:"required"`
	Id        string `json:"id" validate:"required"`
	CreatedAt string `json:"createdAt" validate:"required"`
}
