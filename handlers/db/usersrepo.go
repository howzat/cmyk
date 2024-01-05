package db

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/projects/cmyk-tools/handlers/model"
	"github.com/projects/cmyk-tools/handlers/util"
	"github.com/rs/zerolog"
	"log"
	"strings"
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

	userID := id.String()
	entity := createUserEntity{
		Username:  username,
		Email:     email,
		Pk:        userPk(userID),
		Sk:        userPk(userID),
		Id:        userPk(userID),
		CreatedAt: currentTime.Format(time.RFC3339),
	}

	userEntity, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, err
	}

	entity.Pk = emailPk(email) // allows a uniqueness constraint on email
	entity.Sk = emailPk(email) // allows a uniqueness constraint on email
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

	if err != nil {
		var transactionCancelled *types.TransactionCanceledException
		if errors.As(err, &transactionCancelled) {
			reasons := ExtractCancellationReasons(transactionCancelled.CancellationReasons)
			zerolog.Ctx(ctx).Err(err).Str("reasons", reasons).Msg("failed to persist user")
		} else {
			zerolog.Ctx(ctx).Err(err).Msg("failed to persist user")
		}
	}

	zerolog.Ctx(ctx).Info().Any("res", res).Msg("adding user to table response")

	return entity.ToUser(), err
}

var userPk = func(email string) string { return pk("USER", email) }
var emailPk = func(email string) string { return pk("EMAIL", email) }

func pk(k, v string) string {
	return spew.Sprintf("%s#%s", k, v)
}

func (r *UsersTableRepo) GetUser(ctx context.Context, pk string) (*model.User, error) {
	var entity createUserEntity
	var err error
	if strings.Contains("@", pk) {
		err = r.ddb.GetByKey(ctx, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: userPk(pk)},
			"sk": &types.AttributeValueMemberS{Value: userPk(pk)},
		}, &entity)
	} else {
		err = r.ddb.GetByKey(ctx, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: emailPk(pk)},
			"sk": &types.AttributeValueMemberS{Value: emailPk(pk)},
		}, &entity)
	}

	return entity.ToUser(), err
}

type createUserEntity struct {
	Username  string `dynamodbav:"username" validate:"required"`
	Email     string `dynamodbav:"email" validate:"required"`
	Pk        string `dynamodbav:"pk" validate:"required"`
	Sk        string `dynamodbav:"sk" validate:"required"`
	Id        string `dynamodbav:"id" validate:"required"`
	CreatedAt string `dynamodbav:"createdAt" validate:"required"`
}

func (ue *createUserEntity) ToUser() *model.User {
	return &model.User{
		Username:  ue.Username,
		Email:     ue.Email,
		Id:        ue.Id,
		CreatedAt: ue.CreatedAt,
	}
}
