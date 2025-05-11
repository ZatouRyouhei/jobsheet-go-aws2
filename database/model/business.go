package model

type Business struct {
	ID   int    `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}
