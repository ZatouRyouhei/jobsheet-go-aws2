package client

import (
	"context"
	"jobsheet-go-aws2/database/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Scan() ([]model.Client, error) {
	var clientList []model.Client

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	input := &dynamodb.ScanInput{
		TableName: aws.String("t_client"),
	}

	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &clientList)
	if err != nil {
		return nil, err
	}

	return clientList, nil
}

func GetItem(id int) (*model.Client, error) {
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
		TableName: aws.String("t_client"),
	}
	response, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	targetClient := new(model.Client)
	err = attributevalue.UnmarshalMap(response.Item, targetClient)
	if err != nil {
		return nil, err
	}
	return targetClient, nil
}

func PutItem(targetClient model.Client) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(targetClient)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("t_client"),
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
		TableName: aws.String("t_client"),
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
