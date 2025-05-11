package jobsheet

import (
	"context"
	"jobsheet-go-aws2/database/business_system"
	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/dto"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Scan() ([]model.JobSheet, error) {
	var jobSheetList []model.JobSheet
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	input := &dynamodb.ScanInput{
		TableName: aws.String("t_job_sheet"),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}
	return jobSheetList, nil
}

func GetItem(id string) (*model.JobSheet, error) {
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
		TableName: aws.String("t_job_sheet"),
	}
	response, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	jobSheet := new(model.JobSheet)
	err = attributevalue.UnmarshalMap(response.Item, jobSheet)
	if err != nil {
		return nil, err
	}
	return jobSheet, nil
}

func ScanForIDHeader(idHeader string) ([]model.JobSheet, error) {
	var jobSheetList []model.JobSheet
	filterEx := expression.Name("id").BeginsWith(idHeader)
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
		TableName:                 aws.String("t_job_sheet"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}
	return jobSheetList, nil
}

func ScanForOccur(dateFrom string, dateTo string, systemId int) ([]model.JobSheet, error) {
	var jobSheetList []model.JobSheet
	filterEx := expression.Name("occurdate").GreaterThanEqual(expression.Value(dateFrom))
	filterEx = filterEx.And(expression.Name("occurdate").LessThan(expression.Value(dateTo)))
	filterEx = filterEx.And(expression.Name("business_system_id").Equal(expression.Value(systemId)))
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
		TableName:                 aws.String("t_job_sheet"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}
	return jobSheetList, nil
}

func ScanForComplete(dateFrom string, dateTo string, systemId int) ([]model.JobSheet, error) {
	var jobSheetList []model.JobSheet
	filterEx := expression.Name("completedate").GreaterThanEqual(expression.Value(dateFrom))
	filterEx = filterEx.And(expression.Name("completedate").LessThan(expression.Value(dateTo)))
	filterEx = filterEx.And(expression.Name("business_system_id").Equal(expression.Value(systemId)))
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
		TableName:                 aws.String("t_job_sheet"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}
	return jobSheetList, nil
}

func ScanForUnComplete(systemId int) ([]model.JobSheet, error) {
	var jobSheetList []model.JobSheet
	filterEx := expression.Name("completedate").Equal(expression.Value(""))
	filterEx = filterEx.And(expression.Name("business_system_id").Equal(expression.Value(systemId)))
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
		TableName:                 aws.String("t_job_sheet"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}
	return jobSheetList, nil
}

func Search(condition *dto.RestSearchConditionJobSheet) ([]model.JobSheet, error) {
	var jobSheetList []model.JobSheet
	filterEx := expression.Name("id").NotEqual(expression.Value(0)) // 必ずtrueになる条件
	if condition.Client != 0 {
		filterEx = filterEx.And(expression.Name("client_id").Equal(expression.Value(condition.Client)))
	}
	if condition.Business != 0 {
		businessSystemList, err := business_system.Search(condition.Business)
		if err != nil {
			return nil, err
		}
		for _, businessSystem := range businessSystemList {
			filterEx.And(expression.Name("business_system_id").Equal(expression.Value(businessSystem.ID)))
		}
	}
	if condition.BusinessSystem != 0 {
		filterEx = filterEx.And(expression.Name("business_system_id").Equal(expression.Value(condition.BusinessSystem)))
	}
	if condition.Inquiry != 0 {
		filterEx = filterEx.And(expression.Name("inquiry_id").Equal(expression.Value(condition.Inquiry)))
	}
	if condition.Contact != "" {
		filterEx = filterEx.And(expression.Name("contact_id").Equal(expression.Value(condition.Contact)))
	}
	if condition.Deal != "" {
		filterEx = filterEx.And(expression.Name("deal_id").Equal(expression.Value(condition.Deal)))
	}
	if condition.OccurDateFrom != "" {
		filterEx = filterEx.And(expression.Name("occurdate").GreaterThanEqual(expression.Value(condition.OccurDateFrom)))
	}
	if condition.OccurDateTo != "" {
		filterEx = filterEx.And(expression.Name("occurdate").LessThanEqual(expression.Value(condition.OccurDateTo)))
	}
	if condition.CompleteSign == 1 {
		filterEx = filterEx.And(expression.Name("completedate").NotEqual(expression.Value("")))
	}
	if condition.CompleteSign == 2 {
		filterEx = filterEx.And(expression.Name("completedate").Equal(expression.Value("")))
	}
	if condition.LimitDate != "" {
		filterEx = filterEx.And(expression.Name("limitdate").LessThanEqual(expression.Value(condition.LimitDate)))
	}
	if condition.Keyword != "" {
		reg := "( |　)+"
		keywordArr := regexp.MustCompile(reg).Split(condition.Keyword, -1)
		keywordFilter := expression.Name("id").Equal(expression.Value(0)) // 必ずfalseになる条件
		for _, keyword := range keywordArr {
			keywordFilter = keywordFilter.Or(expression.Name("title").Contains(expression.Value(keyword)))
			keywordFilter = keywordFilter.Or(expression.Name("content").Contains(expression.Value(keyword)))
			keywordFilter = keywordFilter.Or(expression.Name("support").Contains(expression.Value(keyword)))
			keywordFilter = keywordFilter.Or(expression.Name("department").Contains(expression.Value(keyword)))
			keywordFilter = keywordFilter.Or(expression.Name("person").Contains(expression.Value(keyword)))
		}
		filterEx = filterEx.And(keywordFilter)
	}
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
		TableName:                 aws.String("t_job_sheet"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}

	return jobSheetList, nil
}

func SearchForUser(id string) ([]model.JobSheet, error) {
	filterEx := expression.Name("contact_id").Equal(expression.Value(id))
	filterEx = filterEx.Or(expression.Name("deal_id").Equal(expression.Value(id)))
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
		TableName:                 aws.String("t_job_sheet"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	var jobSheetList []model.JobSheet
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}
	return jobSheetList, nil
}

func SearchForBusinessSystem(id int) ([]model.JobSheet, error) {
	filterEx := expression.Name("business_system_id").Equal(expression.Value(id))
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
		TableName:                 aws.String("t_job_sheet"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	var jobSheetList []model.JobSheet
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}
	return jobSheetList, nil
}

func SearchForClient(id int) ([]model.JobSheet, error) {
	filterEx := expression.Name("client_id").Equal(expression.Value(id))
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
		TableName:                 aws.String("t_job_sheet"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	var jobSheetList []model.JobSheet
	err = attributevalue.UnmarshalListOfMaps(response.Items, &jobSheetList)
	if err != nil {
		return nil, err
	}
	return jobSheetList, nil
}

func PutItem(jobSheet model.JobSheet) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(jobSheet)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("t_job_sheet"),
		Item:      item,
	}
	_, err = client.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

func UpdateItem(jobSheet model.JobSheet) error {
	update := expression.Set(expression.Name("completedate"), expression.Value(jobSheet.CompleteDate))
	update.Set(expression.Name("content"), expression.Value(jobSheet.Content))
	update.Set(expression.Name("department"), expression.Value(jobSheet.Department))
	update.Set(expression.Name("limitdate"), expression.Value(jobSheet.LimitDate))
	update.Set(expression.Name("occurdate"), expression.Value(jobSheet.OccurDate))
	update.Set(expression.Name("occurtime"), expression.Value(jobSheet.OccurTime))
	update.Set(expression.Name("person"), expression.Value(jobSheet.Person))
	update.Set(expression.Name("responsetime"), expression.Value(jobSheet.ResponseTime))
	update.Set(expression.Name("support"), expression.Value(jobSheet.Support))
	update.Set(expression.Name("title"), expression.Value(jobSheet.Title))
	update.Set(expression.Name("business_system_id"), expression.Value(jobSheet.BusinessSystemID))
	update.Set(expression.Name("client_id"), expression.Value(jobSheet.ClientID))
	update.Set(expression.Name("contact_id"), expression.Value(jobSheet.ContactID))
	update.Set(expression.Name("deal_id"), expression.Value(jobSheet.DealID))
	update.Set(expression.Name("inquiry_id"), expression.Value(jobSheet.InquiryID))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(cfg)
	paramId, err := attributevalue.Marshal(jobSheet.ID)
	if err != nil {
		return err
	}
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("t_job_sheet"),
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
		TableName: aws.String("t_job_sheet"),
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
