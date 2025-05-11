package attachment

import (
	"context"
	"jobsheet-go-aws2/database/model"
	"mime/multipart"
	"strconv"
	"time"

	"jobsheet-go-aws2/constant"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Scan(id string) ([]model.Attachment, error) {
	var attachmentList []model.Attachment
	filterEx := expression.Name("id").Equal(expression.Value(id))
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
		TableName:                 aws.String("t_attachment"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	response, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &attachmentList)
	if err != nil {
		return nil, err
	}

	return attachmentList, nil
}

func PutItem(attachment model.Attachment, file *multipart.FileHeader) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	// データベースに登録
	client := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(attachment)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("t_attachment"),
		Item:      item,
	}
	_, err = client.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}

	// S3へファイルアップロード
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	ctx := context.Background()
	s3Config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	// ファイル名設定（サブフォルダに入れる場合はスラッシュ区切り。サブフォルダが存在しない場合は勝手に作成される。存在する場合はそのまま利用される。）
	objectKey := attachment.ID + "/" + strconv.Itoa(attachment.SeqNo) + "/" + file.Filename
	s3Client := s3.NewFromConfig(s3Config)
	s3Input := &s3.PutObjectInput{
		Bucket: aws.String(constant.BUCKET_NAME),
		Key:    aws.String(objectKey),
		Body:   src,
	}
	// S3へのアップロード処理
	_, err = s3Client.PutObject(ctx, s3Input)
	if err != nil {
		return err
	}
	err = s3.NewObjectExistsWaiter(s3Client).Wait(ctx, &s3.HeadObjectInput{Bucket: aws.String(constant.BUCKET_NAME), Key: aws.String(objectKey)}, time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func GetItem(id string, seqNo int) (*model.Attachment, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	paramId, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}
	paramSeqNo, err := attributevalue.Marshal(seqNo)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id":    paramId,
			"seqno": paramSeqNo,
		},
		TableName: aws.String("t_attachment"),
	}
	response, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	attachment := new(model.Attachment)
	err = attributevalue.UnmarshalMap(response.Item, attachment)
	if err != nil {
		return nil, err
	}
	return attachment, nil
}

// 特定の業務日誌に関連する添付ファイルを全て削除する。
// idのみを指定して削除する。
func DeleteItems(id string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	deleteAttachmentList, err := Scan(id)
	if err != nil {
		return err
	}
	// S3から削除
	s3Client := s3.NewFromConfig(cfg)
	for _, targetAttachment := range deleteAttachmentList {
		objectKey := id + "/" + strconv.Itoa(targetAttachment.SeqNo) + "/" + targetAttachment.FileName
		input := &s3.DeleteObjectInput{
			Bucket: aws.String(constant.BUCKET_NAME),
			Key:    aws.String(objectKey),
		}
		_, err = s3Client.DeleteObject(context.TODO(), input)
		if err != nil {
			return err
		} else {
			err = s3.NewObjectNotExistsWaiter(s3Client).Wait(
				context.TODO(), &s3.HeadObjectInput{Bucket: aws.String(constant.BUCKET_NAME), Key: aws.String(objectKey)}, time.Minute)
			if err != nil {
				return err
			}
		}
	}
	// DBから削除
	client := dynamodb.NewFromConfig(cfg)
	paramId, err := attributevalue.Marshal(id)
	if err != nil {
		return err
	}
	for _, targetAttachment := range deleteAttachmentList {
		paramSeqNo, err := attributevalue.Marshal(targetAttachment.SeqNo)
		if err != nil {
			return err
		}
		// パーティションキーとソートキーを両方指定する必要がある。
		input := &dynamodb.DeleteItemInput{
			TableName: aws.String("t_attachment"),
			Key: map[string]types.AttributeValue{
				"id":    paramId,
				"seqno": paramSeqNo,
			},
		}
		_, err = client.DeleteItem(context.TODO(), input)
		if err != nil {
			return err
		}
	}
	return nil
}

// 特定の添付ファイルを削除する
func DeleteItem(id string, seqNo int) error {
	targetAttachment, err := GetItem(id, seqNo)
	if err != nil {
		return err
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	s3Client := s3.NewFromConfig(cfg)
	// 削除対象ファイル指定（ファイルを削除した結果フォルダの中身が空になるとフォルダも自動的に削除される）
	objectKey := id + "/" + strconv.Itoa(seqNo) + "/" + targetAttachment.FileName
	s3Input := &s3.DeleteObjectInput{
		Bucket: aws.String(constant.BUCKET_NAME),
		Key:    aws.String(objectKey),
	}
	_, err = s3Client.DeleteObject(context.TODO(), s3Input)
	if err != nil {
		return err
	} else {
		err = s3.NewObjectNotExistsWaiter(s3Client).Wait(
			context.TODO(), &s3.HeadObjectInput{Bucket: aws.String(constant.BUCKET_NAME), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			return err
		}
	}
	// DBから削除
	client := dynamodb.NewFromConfig(cfg)
	paramId, err := attributevalue.Marshal(id)
	if err != nil {
		return err
	}
	paramSeqNo, err := attributevalue.Marshal(seqNo)
	if err != nil {
		return err
	}
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("t_attachment"),
		Key: map[string]types.AttributeValue{
			"id":    paramId,
			"seqno": paramSeqNo,
		},
	}
	_, err = client.DeleteItem(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
