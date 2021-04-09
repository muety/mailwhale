package dto

import "github.com/muety/mailwhale/types"

type UserUpdate struct {
	Password string              `json:"password"`
	Senders  types.MailAddresses `json:"senders"`
}

func (u *UserUpdate) GetSenders(user *types.User) []types.SenderAddress {
	senders := make([]types.SenderAddress, len(u.Senders))
	for i := range u.Senders {
		senders[i] = types.SenderAddress{
			MailAddress: u.Senders[i],
			Verified:    user.HasVerifiedSender(u.Senders[i]),
		}
	}
	return senders
}
