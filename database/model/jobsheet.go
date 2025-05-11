package model

type JobSheet struct {
	ID               string  `dynamodbav:"id"`
	CompleteDate     string  `dynamodbav:"completedate"`
	Content          string  `dynamodbav:"content"`
	Department       string  `dynamodbav:"department"`
	LimitDate        string  `dynamodbav:"limitdate"`
	OccurDate        string  `dynamodbav:"occurdate"`
	OccurTime        string  `dynamodbav:"occurtime"`
	Person           string  `dynamodbav:"person"`
	ResponseTime     float64 `dynamodbav:"responsetime"`
	Support          string  `dynamodbav:"support"`
	Title            string  `dynamodbav:"title"`
	BusinessSystemID int     `dynamodbav:"business_system_id"`
	ClientID         int     `dynamodbav:"client_id"`
	ContactID        string  `dynamodbav:"contact_id"`
	DealID           string  `dynamodbav:"deal_id"`
	InquiryID        int     `dynamodbav:"inquiry_id"`
}
