package model


type User struct{
	Email string `json:"email"`
	PasswordHash string `json:"passwordhash"`
	Address string 	`json:"address"`
}