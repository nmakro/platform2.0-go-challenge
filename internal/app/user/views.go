package user

import "encoding/json"

type PublicUserView struct {
	user *User
}

func (p *PublicUserView) MarshalJSON() ([]byte, error) {
	type publicUser struct {
		*User
		Password  *string `json:"password,omitempty"`
		FirstName *string `json:"first_name,omitempty"`
		LastName  *string `json:"last_name,omitempty"`
	}

	pu := publicUser{
		User: p.user,
	}

	return json.Marshal(&pu)
}
