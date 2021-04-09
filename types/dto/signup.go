package dto

type Signup struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
