package model

type Attachment struct {
	ID       string `dynamodbav:"id"`
	SeqNo    int    `dynamodbav:"seqno"`
	FileName string `dynamodbav:"filename"`
}
