package business_system

import (
	"context"
	"jobsheet-go-aws2/database/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Scan() ([]model.BusinessSystem, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	input := &dynamodb.ScanInput{
		TableName: aws.String("t_business_system"),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	var businessSystemList []model.BusinessSystem
	err = attributevalue.UnmarshalListOfMaps(response.Items, &businessSystemList)
	if err != nil {
		return nil, err
	}
	return businessSystemList, nil
}

func GetItem(id int) (*model.BusinessSystem, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	paramId, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": paramId,
		},
		TableName: aws.String("t_business_system"),
	}
	response, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	businessSystem := new(model.BusinessSystem)
	err = attributevalue.UnmarshalMap(response.Item, businessSystem)
	if err != nil {
		return nil, err
	}
	return businessSystem, nil
}

func Search(id int) ([]model.BusinessSystem, error) {
	var businessSystemList []model.BusinessSystem
	filterEx := expression.Name("business_id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithFilter(filterEx).Build()
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	input := &dynamodb.ScanInput{
		TableName:                 aws.String("t_business_system"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &businessSystemList)
	if err != nil {
		return nil, err
	}
	return businessSystemList, nil
}

func PutItem(businessSystem model.BusinessSystem) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(businessSystem)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("t_business_system"),
		Item:      item,
	}
	_, err = client.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItem(id int) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(cfg)
	paramId, err := attributevalue.Marshal(id)
	if err != nil {
		return err
	}
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("t_business_system"),
		Key: map[string]types.AttributeValue{
			"id": paramId,
		},
	}
	_, err = client.DeleteItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
