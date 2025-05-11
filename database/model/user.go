package model

type User struct {
	Id       string `dynamodbav:"id"`
	Password string `dynamodbav:"password"`
	Name     string `dynamodbav:"name"`
	SeqNo    int    `dynamodbav:"seqno"`
}
