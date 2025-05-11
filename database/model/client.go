package model

type Client struct {
	ID   int    `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}
