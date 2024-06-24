package app

type userData struct {
	id                int
	firstName         string
	lastName          string
	license           string
	email             string
	phone             string
	encryptedPassword string
	isActive          bool
	isAdmin           bool
	notifications     bool
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type GetUserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	License   string `json:"license"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"isAdmin"`
}

type UserRegisterRequest struct {
	License       string `json:"license"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Notifications bool   `json:"notifications"`
}

type UserRegisterResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`

	ID        int    `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	License   string `json:"license"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"isAdmin"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`

	ID        int    `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	License   string `json:"license"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"isAdmin"`
}
