package model

type User struct {
	Username  string `dynamodbav:"username" json:"username" validate:"required"`
	Email     string `dynamodbav:"email" json:"email" validate:"required"`
	Id        string `dynamodbav:"id" json:"id" validate:"required"`
	CreatedAt string `dynamodbav:"createdAt" json:"createdAt" validate:"required"`
}
