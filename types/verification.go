package types

import "github.com/google/uuid"

const (
	VerificationScopeUser   = "scope_user"
	VerificationScopeSender = "scope_sender_address"
)

type Verification struct {
	Token   string `boltholdKey:"Token"`
	UserId  string
	Scope   string
	Subject string
}

func NewVerification(user *User, scope, subject string) *Verification {
	return &Verification{
		Token:   uuid.New().String(),
		UserId:  user.ID,
		Scope:   scope,
		Subject: subject,
	}
}
