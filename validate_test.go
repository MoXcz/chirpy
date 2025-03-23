package main

import "testing"

func Test_removeProfanity(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		profaneBody string
		want        string
	}{
		{
			name:        "Remove profanity: kerfuffle",
			profaneBody: "This is a kerfuffle opinion I need to share with the world",
			want:        "This is a **** opinion I need to share with the world",
		},
		{
			name:        "Remove profanity: sharbert",
			profaneBody: "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			want:        "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			name:        "Remove profanity: all",
			profaneBody: "I really need a kerfuffle to go to bed sooner, Fornax !",
			want:        "I really need a **** to go to bed sooner, **** !",
		},
		{
			name:        "Do not remove profanity: punctuation",
			profaneBody: "fornax!",
			want:        "fornax!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeProfanity(tt.profaneBody)
			if got != tt.want {
				t.Errorf("removeProfanity() = %v, want %v", got, tt.want)
			}
		})
	}
}
