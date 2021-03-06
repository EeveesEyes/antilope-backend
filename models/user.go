package models

type User struct {
	Id       int
	Username string
	Email    string
	// password
	Hash       string
	PepperID   int
	CurrentJWT string
}

type CreateUserData struct {
	Username string
	Email    string
	Password string
}

type UserRequest struct {
	UserData CreateUserData `json:"user"`
}

func NewUser(id int, username, email, hash string, pepperID int) *User {
	return &User{
		Id:       id,
		Username: username,
		Email:    email,
		Hash:     hash,
		PepperID: pepperID,
	}
}
