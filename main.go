package main

import (
	"context"
	"jobsheet-go-aws2/logger"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	// システムログ設定
	logger.LogInit()

	// echoを起動
	e := echo.New()

	// ルーティング設定
	SetRoute(e)

	// echoのログ取得
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: os.Stdout,
	}))

	echoLambda = echoadapter.New(e)

	lambda.Start(Handler)
}
