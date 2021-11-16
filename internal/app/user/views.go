package user

import "encoding/json"

func (p *User) MarshalJSON() ([]byte, error) {
	type PublicUserAlias User

	return json.Marshal(&struct {
		FirstName *struct{} `json:"first_name,omitempty"`
		LastName  *struct{} `json:"last_name,omitempty"`
		*PublicUserAlias
	}{
		PublicUserAlias: (*PublicUserAlias)(p),
	})
}
