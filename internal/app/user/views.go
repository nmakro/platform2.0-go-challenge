package user

import "encoding/json"

func (p *User) MarshalJSON() ([]byte, error) {
	type PublicUserAlias User

	return json.Marshal(&struct {
		UserID    *struct{} `json:"user_id,omitempty"`
		FirstName *struct{} `json:"first_name,omitempty"`
		LastName  *struct{} `json:"last_name,omitempty"`
		Password  *struct{} `json:"password"`
		*PublicUserAlias
	}{
		PublicUserAlias: (*PublicUserAlias)(p),
	})
}
