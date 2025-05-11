package model

type Holiday struct {
	Holiday string `dynamodbav:"holiday"`
	Name    string `dynamodbav:"name"`
}
