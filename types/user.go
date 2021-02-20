package types

import "github.com/muety/mailwhale/util"

type User struct {
	ID       string `json:"id"`
	Password string `json:"-"`
}

type Signup struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *User) IsValid() bool {
	return util.IsEmail(u.ID) && len(u.Password) > 0
}
