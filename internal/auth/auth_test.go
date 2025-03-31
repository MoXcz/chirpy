package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
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
			name:          "correct password",
			password:      "thisisaverycomplexpasswordindeed",
			wantErr:       false,
			checkPassword: "thisisaverycomplexpasswordindeed",
		},
		{
			name:          "incorrect password",
			password:      "onewouldargueiflongpasswordsarereallysecure",
			wantErr:       true,
			checkPassword: "pepito2721",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cHash, err := HashPassword(tt.password)
			if err != nil {
				fmt.Println(err)
			}
			gotErr := CheckPasswordHash(cHash, tt.checkPassword)
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

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			want:        userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			want:        uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			want:        uuid.Nil,
			wantErr:     true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, gotErr := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
			if gotUserID != tt.want {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.want)
			}
		})
	}
}
