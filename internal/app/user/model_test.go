package user

import "testing"

func TestValidateUser(t *testing.T) {

	tests := []struct {
		name   string
		args   AddUserCommand
		errMSg string
	}{
		{
			name: "valid",
			args: AddUserCommand{
				UserName: "test",
				Email:    "test@host.com",
				Password: "@3das@!DZ",
			},
		},
		{
			name: "small password",
			args: AddUserCommand{
				UserName: "test",
				Email:    "test@host.com",
				Password: "!@3s",
			},
			errMSg: "password must be at least eight characters long",
		},
		{
			name: "no digits password",
			args: AddUserCommand{
				UserName: "test",
				Email:    "test@host.com",
				Password: "asdfkkawwf",
			},
			errMSg: "password must have at least one digit",
		},
		{
			name: "no valid email",
			args: AddUserCommand{
				UserName: "test",
				Email:    "test.host",
				Password: "@3das@!DZ",
			},
			errMSg: "user email test.host is not valid",
		},
		{
			name: "no valid email",
			args: AddUserCommand{
				FirstName: "witdDigits1",
				UserName:  "test",
				Email:     "test@host.com",
				Password:  "@3das@!DZ",
			},
			errMSg: "only letters are allowed in first and last name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateUser(tt.args); err != nil {
				if err.Error() != tt.errMSg {
					t.Errorf("ValidateUser() error = %s, wantErr %s", err.Error(), tt.errMSg)
				}
			}
		})
	}
}
