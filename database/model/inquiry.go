package model

type Inquiry struct {
	ID   int    `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}
