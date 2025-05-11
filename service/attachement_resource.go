package service

import (
	"context"
	"io"
	"jobsheet-go-aws2/constant"
	"jobsheet-go-aws2/database/attachment"
	"jobsheet-go-aws2/database/model"
	"jobsheet-go-aws2/dto"
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func RegistAttachment(c echo.Context) error {
	id := c.Param("id")
	form, err := c.MultipartForm()
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	files := form.File["file"]
	for _, file := range files {
		nextSeqNo := 1
		attachmentList, err := attachment.Scan(id)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
		if len(attachmentList) > 0 {
			// 連番順の降順にソートする
			sort.Slice(attachmentList, func(i, j int) bool {
				return attachmentList[i].SeqNo > attachmentList[j].SeqNo
			})
			nextSeqNo = attachmentList[0].SeqNo + 1
		}
		targetAttachement := new(model.Attachment)
		targetAttachement.ID = id
		targetAttachement.SeqNo = nextSeqNo
		targetAttachement.FileName = file.Filename
		err = attachment.PutItem(*targetAttachement, file)
		if err != nil {
			slog.Error("Error", slog.Any("error", err))
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	// 返信用の添付ファイルリストを作成
	respAttachmentList, err := attachment.Scan(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var restAttachmentList []dto.RestAttachment
	for _, attachment := range respAttachmentList {
		restAttachmentList = append(restAttachmentList, dto.RestAttachment{
			SeqNo:    attachment.SeqNo,
			FileName: attachment.FileName,
		})
	}
	// 連番の昇順にソート
	sort.Slice(restAttachmentList, func(i, j int) bool {
		return restAttachmentList[i].SeqNo < restAttachmentList[j].SeqNo
	})
	return c.JSON(http.StatusCreated, restAttachmentList)
}

func DownloadAttachment(c echo.Context) error {
	id := c.Param("id")
	seqNoStr := c.Param("seqNo")
	seqNo, err := strconv.Atoi(seqNoStr)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	targetAttachment, err := attachment.GetItem(id, seqNo)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	// S3からダウンロードする
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	// ファイル名指定
	objectKey := id + "/" + seqNoStr + "/" + targetAttachment.FileName
	// S3からファイル取得処理
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(constant.BUCKET_NAME),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer result.Body.Close()
	targetFile, err := io.ReadAll(result.Body)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	response := c.Response()
	response.Writer.Header().Set("Content-Disposition", "attachment; filename="+targetAttachment.FileName)
	return c.Blob(http.StatusOK, "application/octet-stream", targetFile)
}

func DeleteAttachment(c echo.Context) error {
	id := c.Param("id")
	seqNoStr := c.Param("seqNo")
	seqNo, err := strconv.Atoi(seqNoStr)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = attachment.DeleteItem(id, seqNo)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	// 返信用の添付ファイルリストを作成
	respAttachmentList, err := attachment.Scan(id)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	var restAttachmentList []dto.RestAttachment
	for _, attachment := range respAttachmentList {
		restAttachmentList = append(restAttachmentList, dto.RestAttachment{
			SeqNo:    attachment.SeqNo,
			FileName: attachment.FileName,
		})
	}
	// 連番の昇順にソート
	sort.Slice(restAttachmentList, func(i, j int) bool {
		return restAttachmentList[i].SeqNo < restAttachmentList[j].SeqNo
	})
	return c.JSON(http.StatusCreated, restAttachmentList)
}
