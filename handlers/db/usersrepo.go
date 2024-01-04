package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/projects/cmyk-tools/handlers/model"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/rs/zerolog"
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

func (r *UsersTableRepo) AddUser(ctx context.Context, username, email string) (*model.User, error) {

	currentTime, id, err := util.CurrentTimeAndULID(r.clock)
	if err != nil {
		log.Fatal(err)
	}

	entity := createUserEntity{
		Username:  username,
		Email:     email,
		Pk:        id.String(),
		Id:        id.String(),
		CreatedAt: currentTime.Format(time.RFC3339),
	}

	userEntity, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, err
	}

	entity.Pk = email // allows a uniqueness constraint on email
	emailEntity, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, err
	}

	res, err := r.ddb.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					Item:                userEntity,
					TableName:           aws.String(r.ddb.Tablename),
					ConditionExpression: aws.String("attribute_not_exists(pk)"),
				},
			},
			{
				Put: &types.Put{
					Item:                emailEntity,
					TableName:           aws.String(r.ddb.Tablename),
					ConditionExpression: aws.String("attribute_not_exists(pk)"),
				},
			},
		}})

	zerolog.Ctx(ctx).Info().Any("res", res).Msg("adding user to table response")

	return entity.ToUser(), err
}

func (r *UsersTableRepo) GetUser(ctx context.Context, pk string) (*model.User, error) {
	var entity createUserEntity
	err := r.ddb.GetByKey(ctx, map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
	}, &entity)

	return entity.ToUser(), err
}

type createUserEntity struct {
	Username  string `dynamodbav:"username" validate:"required"`
	Email     string `dynamodbav:"email" validate:"required"`
	Pk        string `dynamodbav:"pk" validate:"required"`
	Id        string `dynamodbav:"id" validate:"required"`
	CreatedAt string `dynamodbav:"createdAt" validate:"required"`
}

func (ue *createUserEntity) ToUser() *model.User {
	return &model.User{
		Username:  ue.Username,
		Email:     ue.Email,
		Pk:        ue.Pk,
		Id:        ue.Id,
		CreatedAt: ue.CreatedAt,
	}
}
