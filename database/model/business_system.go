package model

type BusinessSystem struct {
	ID         int    `dynamodbav:"id"`
	Name       string `dynamodbav:"name"`
	BusinessID int    `dynamodbav:"business_id"`
}
