package auth_test

import (
	"fmt"
	"testing"

	"github.com/MoXcz/chirpy/internal/auth"
)

func TestCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		hash          string
		password      string
		checkPassword string
		wantErr       bool
	}{
		{
			password:      "thisisaverycomplexpasswordindeed",
			wantErr:       false,
			checkPassword: "thisisaverycomplexpasswordindeed",
		},
		{
			password:      "onewouldargueiflongpasswordsarereallysecure",
			wantErr:       true,
			checkPassword: "pepito2721",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cHash, err := auth.HashPassword(tt.password)
			if err != nil {
				fmt.Println(err)
			}
			gotErr := auth.CheckPasswordHash(cHash, tt.checkPassword)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CheckPasswordHash() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CheckPasswordHash() succeeded unexpectedly")
			}
		})
	}
}
