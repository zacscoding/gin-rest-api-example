package model

import (
	"fmt"
	"testing"
)

func TestAccountString(t *testing.T) {
	cases := []struct {
		Name string
		Acc  Account
	}{
		{
			Name: "empty password",
			Acc: Account{
				ID:       1,
				Username: "user1",
				Email:    "user1@email.com",
				Password: "",
			},
		}, {
			Name: "exist password",
			Acc: Account{
				ID:       1,
				Username: "user1",
				Email:    "user1@email.com",
				Password: "pass1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			expected := fmt.Sprintf("{%d %s %s %s}", tc.Acc.ID, tc.Acc.Username, tc.Acc.Email, "[PROTECTED]")
			if got := tc.Acc.String(); expected != got {
				t.Errorf("String() wanted:%s, got:%s", expected, got)
			}
		})
	}
}
