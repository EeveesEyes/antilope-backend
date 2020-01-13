package models

type User struct {
	Id       string
	Username string
	Email    string
	// password
	Hash     string
	PepperID int
}

type CreateUserData struct {
	Username string
	Email    string
	Password string
}

type UserRequest struct {
	UserData CreateUserData `json:"user"`
}

type UserResponse struct {
	UserData UserAuthResponse `json:"user"`
}

type UserAuthResponse struct {
	Username string
	Email    string
	Token    string
}

func NewUser(username, email, hash string, pepperID int) *User {
	return &User{
		Username: username,
		Email:    email,
		Hash:     hash,
		PepperID: pepperID,
	}
}
