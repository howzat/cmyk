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
	"time"
)

type UsersRepo struct {
	ddb   DynamoRepository
	clock util.Clock
}

func NewUsersTableRepo(ctx context.Context, region string) (*UsersRepo, error) {
	instance, err := NewInstance(ctx, region, UsersTableEnvKey)
	if err != nil {
		return nil, err
	}

	return &UsersRepo{
		ddb:   *instance,
		clock: util.NewRealClock(),
	}, nil
}

func (r *UsersRepo) AddTestUser(ctx context.Context, user model.User, lifespan util.Lifespan) (*model.User, error) {
	ttlExpiry := util.TestLifespan(lifespan, time.Now())
	return r.addUser(ctx, user, &ttlExpiry)
}

func (r *UsersRepo) AddUser(ctx context.Context, user model.User) (*model.User, error) {
	return r.addUser(ctx, user, nil)
}
func (r *UsersRepo) addUser(ctx context.Context, user model.User, ttl *int64) (*model.User, error) {

	entity := createUserEntity(user, ttl)
	userEntity, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, err
	}

	// email entity sets pk and sk to the email which ensures uniqueness
	// (needs to be region pinned to avoid consistency races)
	emailEntity, err := attributevalue.MarshalMap(createEmailUniquenessEntity(user, ttl))
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

	if err != nil {
		return nil, err
	}

	return entity.ToUser()
}

func createEmailUniquenessEntity(user model.User, ttl *int64) emailUniquenessEntity {

	entity := emailUniquenessEntity{
		Pk: emailPk(user.Email),
		Sk: emailPk(user.Email),
	}

	if ttl != nil && *ttl > 0 {
		entity.ExpireAt = *ttl
	}

	return entity
}

func createUserEntity(user model.User, ttl *int64) userEntity {

	entity := userEntity{
		Email:     user.Email,
		Pk:        usernamePK(user.Id),
		Sk:        usernamePK(user.Id),
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	if ttl != nil && *ttl > 0 {
		entity.ExpireAt = *ttl
	}

	return entity
}

var usernamePK = func(email string) string { return pk("USERNAME", email) }
var emailPk = func(email string) string { return pk("USEREMAIL", email) }

func pk(k, v string) string {
	return spew.Sprintf("%s#%s", k, v)
}

func (r *UsersRepo) GetUserByID(ctx context.Context, pk string) (*model.User, error) {

	var entity userEntity
	err := r.ddb.GetByKey(ctx, map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: usernamePK(pk)},
		"sk": &types.AttributeValueMemberS{Value: usernamePK(pk)},
	}, &entity)

	if err != nil {
		return nil, err
	}

	return entity.ToUser()
}

type userEntity struct {
	Pk        string `dynamodbav:"pk" validate:"required"`
	Sk        string `dynamodbav:"sk" validate:"required"`
	CreatedAt string `dynamodbav:"createdAt" validate:"required"`
	Email     string `dynamodbav:"email" validate:"required"`
	Name      string `dynamodbav:"name" validate:"required"`
	ExpireAt  int64  `dynamodbav:"ttl"`
}

func (ue *userEntity) ToUser() (*model.User, error) {
	timestamp, err := time.Parse(time.RFC3339, ue.CreatedAt)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Id:        ue.Pk,
		Name:      ue.Name,
		Email:     ue.Email,
		CreatedAt: timestamp,
	}

	if ue.ExpireAt > 0 {
		user.Ttl = &ue.ExpireAt
	}

	return &user, nil
}

type emailUniquenessEntity struct {
	Pk       string `dynamodbav:"pk" validate:"required"`
	Sk       string `dynamodbav:"sk" validate:"required"`
	ExpireAt int64  `dynamodbav:"ttl"`
}
