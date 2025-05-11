package user

import (
	"context"
	"errors"
	"jobsheet-go-aws2/database/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Scan() ([]model.User, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	input := &dynamodb.ScanInput{
		TableName: aws.String("t_user"),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	var userList []model.User
	err = attributevalue.UnmarshalListOfMaps(response.Items, &userList)
	if err != nil {
		return nil, err
	}

	return userList, nil
}

func PassCheck(id string, password string) (*model.User, error) {
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
		TableName: aws.String("t_user"),
	}
	response, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	user := new(model.User)
	err = attributevalue.UnmarshalMap(response.Item, user)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func GetItem(id string) (*model.User, error) {
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
		TableName: aws.String("t_user"),
	}
	response, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	user := new(model.User)
	err = attributevalue.UnmarshalMap(response.Item, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func PutItem(user model.User) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("t_user"),
		Item:      item,
	}
	_, err = client.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

func UpdateItem(user model.User) error {
	update := expression.Set(expression.Name("name"), expression.Value(user.Name))
	if user.Password != "" {
		update.Set(expression.Name("password"), expression.Value(user.Password))
	}
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(cfg)
	paramId, err := attributevalue.Marshal(user.Id)
	if err != nil {
		return err
	}
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("t_user"),
		Key: map[string]types.AttributeValue{
			"id": paramId,
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}
	_, err = client.UpdateItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSeqNo(user model.User) error {
	update := expression.Set(expression.Name("seqno"), expression.Value(user.SeqNo))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(cfg)
	paramId, err := attributevalue.Marshal(user.Id)
	if err != nil {
		return err
	}
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("t_user"),
		Key: map[string]types.AttributeValue{
			"id": paramId,
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}
	_, err = client.UpdateItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItem(id string) error {
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
		TableName: aws.String("t_user"),
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
