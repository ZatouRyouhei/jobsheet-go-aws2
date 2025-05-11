package holiday

import (
	"context"
	"jobsheet-go-aws2/database/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Scan() ([]model.Holiday, error) {
	var holidayList []model.Holiday

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	input := &dynamodb.ScanInput{
		TableName: aws.String("t_holiday"),
	}

	// scanPaginator := dynamodb.NewScanPaginator(client, input)
	// for scanPaginator.HasMorePages() {
	// 	response, err := scanPaginator.NextPage(context.TODO())
	// 	if err != nil {
	// 		return c.String(http.StatusBadRequest, "bad request")
	// 	}
	// 	var holidayPage []model.Holiday
	// 	err = attributevalue.UnmarshalListOfMaps(response.Items, &holidayPage)
	// 	if err != nil {
	// 		return c.String(http.StatusBadRequest, "bad request")
	// 	}
	// 	holidayList = append(holidayList, holidayPage...)
	// }

	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &holidayList)
	if err != nil {
		return nil, err
	}

	return holidayList, nil
}

func BatchWriteItem(holidays []model.Holiday) (int, error) {
	var err error
	var item map[string]types.AttributeValue
	written := 0
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return 0, err
	}
	client := dynamodb.NewFromConfig(cfg)
	for start < len(holidays) {
		var writeReqs []types.WriteRequest
		if end > len(holidays) {
			end = len(holidays)
		}
		for _, holiday := range holidays[start:end] {
			item, err = attributevalue.MarshalMap(holiday)
			if err != nil {
				return 0, err
			} else {
				writeReqs = append(writeReqs, types.WriteRequest{
					PutRequest: &types.PutRequest{
						Item: item,
					},
				})
			}
		}
		_, err = client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				"t_holiday": writeReqs,
			},
		})
		if err != nil {
			return 0, err
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}
	return written, err
}
