package models

type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FName    string `json:"fname"`
	LName    string `json:"lname"`
}

type TokenData struct {
	Token string `json:"token"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CheckAuthData struct {
	Token string `json:"token"`
}
