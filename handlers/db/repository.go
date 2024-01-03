package db

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const UsersTableEnvKey = "DYNAMODB_USERS_TABLE"

type Transaction interface {
	TransactPut(ctx *context.Context, items []*types.TransactWriteItem) error
	TransactWriteItem(ctx context.Context, item interface{}) (*types.TransactWriteItem, error)
}

type DynamoRepository struct {
	Tablename string
	Client    *dynamodb.Client
}

func NewInstance(ctx context.Context, regionKey string, tablenameKey string) DynamoRepository {
	region := os.Getenv(regionKey)
	tablename := os.Getenv(tablenameKey)

	return NewInstanceWithValues(*zerolog.Ctx(ctx), region, tablename)
}

// NewInstanceWithValues takes region and tablename as values instead of environment variable keys to be looked up.
// It allows callers to take responsibility for env var lookup which can make it easier to tell which env vars a lambda needs.
func NewInstanceWithValues(log zerolog.Logger, region string, tablename string) DynamoRepository {
	log.With().Str("region", region).Str("tablename", tablename).Logger()
	log.Debug().Msg("DynamoDB instance configuration")

	return DynamoRepository{
		Client:    NewDynamoDB(region),
		Tablename: tablename,
	}
}

func (r *DynamoRepository) GetTablename() string {
	return r.Tablename
}

func (r *DynamoRepository) GetClient() *dynamodb.Client {
	return r.Client
}

func NewDynamoDB(region string) *dynamodb.Client {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		log.Fatal(err)
	}

	return dynamodb.NewFromConfig(awsConfig)
}

type NotFoundError struct {
	StatusCode int
	Err        error
}

func NewNotFoundError(err error) NotFoundError {
	return NotFoundError{
		StatusCode: 404,
		Err:        err,
	}
}
func (m NotFoundError) Error() string {
	return m.Err.Error()
}

func (r *DynamoRepository) GetByKey(ctx context.Context, key map[string]types.AttributeValue, model interface{}) error {

	result, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.Tablename),
		Key:       key,
	})

	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to get item from dynamodb")
		return err
	}
	if result.Item == nil {
		return NewNotFoundError(errors.New(spew.Sprintf("item not found for key [%s]", key)))
	}

	err = attributevalue.UnmarshalMap(result.Item, model)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("item", result.Item).Msg("Failed to unmarshal record item")
		return err
	}
	return nil
}

func (r *DynamoRepository) Put(ctx context.Context, model interface{}) error {

	av, err := attributevalue.MarshalMap(model)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to marshal into dynamodb map")
		return err
	}
	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.Tablename),
	})
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to persist model")
		return err
	}

	zerolog.Ctx(ctx).Debug().Any("model", model).Msg("Persisted")

	return nil
}

func (r *DynamoRepository) TransactPut(ctx context.Context, items []types.TransactWriteItem) error {

	_, err := r.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	})

	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("items", items).Msg("Failed to persist all transactional writes")
		return err
	}
	zerolog.Ctx(ctx).Debug().Any("items", items).Msg("Persisted")
	return nil
}

func (r *DynamoRepository) TransactWriteItem(ctx context.Context, item interface{}) (*types.TransactWriteItem, error) {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to marshal transaction item into dynamodb map")
		return nil, err
	}
	return &types.TransactWriteItem{
		Put: &types.Put{
			TableName: aws.String(r.Tablename),
			Item:      av,
		},
	}, nil
}

func (r *DynamoRepository) Update(ctx context.Context, input *dynamodb.UpdateItemInput) error {
	input.TableName = aws.String(r.Tablename)
	_, err := r.Client.UpdateItem(ctx, input)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to update model")
		return err
	}
	return nil
}

func (r *DynamoRepository) Scan(ctx context.Context, models interface{}) error {
	result, err := r.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.Tablename),
	})
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to scan table")
		return err
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, models)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to unmarshal list of items")
		return err
	}
	return nil
}

func (r *DynamoRepository) Query(ctx context.Context, input *dynamodb.QueryInput, models interface{}) error {
	input.TableName = aws.String(r.Tablename)
	result, err := r.Client.Query(ctx, input)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to query table")
		return err
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, models)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to unmarshal list of items")
		return err
	}
	return nil
}

func (r *DynamoRepository) Delete(ctx context.Context, ID string) error {
	input := &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: ID}},
		TableName: aws.String(r.Tablename),
	}
	if _, err := r.Client.DeleteItem(ctx, input); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to delete item")
		return err
	}

	return nil
}

type MultiWriteItem struct {
	TableName string
	Model     interface{}
}

func (r *DynamoRepository) TransactPutMultiTable(ctx context.Context, arr []MultiWriteItem) error {
	items := make([]types.TransactWriteItem, 0)

	for _, v := range arr {
		av, err := attributevalue.MarshalMap(v.Model)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to marshal into dynamodb map")
			return err
		}
		items = append(items, types.TransactWriteItem{
			Put: &types.Put{
				Item:      av,
				TableName: aws.String(v.TableName),
			},
		})
	}
	_, err := r.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	})
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to persist all write items in transaction")
		return err
	}
	zerolog.Ctx(ctx).Debug().Any("items", items).Msg("Persisted transactional items")
	return nil
}
