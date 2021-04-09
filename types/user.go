package types

import "github.com/muety/mailwhale/util"

type User struct {
	ID       string          `json:"id" boltholdKey:"ID"`
	Password string          `json:"-"`
	Senders  []SenderAddress `json:"senders"`
}

type Signup struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *User) IsValid() bool {
	return util.IsEmail(u.ID) && len(u.Password) > 0
}

func (u *User) HasSender(sender MailAddress) bool {
	return u.findSender(sender) != nil
}

func (u *User) HasVerifiedSender(sender MailAddress) bool {
	s := u.findSender(sender)
	return s != nil && s.Verified
}

func (u *User) Sanitize() *User {
	if u.Senders == nil {
		u.Senders = []SenderAddress{}
	}
	return u
}

func (u *User) findSender(sender MailAddress) *SenderAddress {
	for _, v := range u.Senders {
		if v.MailAddress == sender {
			return &v
		}
	}
	return nil
}
