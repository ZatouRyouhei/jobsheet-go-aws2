package inquiry

import (
	"context"
	"jobsheet-go-aws2/database/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Scan() ([]model.Inquiry, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	input := &dynamodb.ScanInput{
		TableName: aws.String("t_inquiry"),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	var inquiryList []model.Inquiry
	err = attributevalue.UnmarshalListOfMaps(response.Items, &inquiryList)
	if err != nil {
		return nil, err
	}

	return inquiryList, nil
}

func GetItem(id int) (*model.Inquiry, error) {
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
		TableName: aws.String("t_inquiry"),
	}
	response, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	inquiry := new(model.Inquiry)
	err = attributevalue.UnmarshalMap(response.Item, inquiry)
	if err != nil {
		return nil, err
	}
	return inquiry, nil
}
