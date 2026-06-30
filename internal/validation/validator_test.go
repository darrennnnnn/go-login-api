package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestPasswordValidation(t *testing.T) {
	v := validator.New()
	if err := v.RegisterValidation("password", Password); err != nil {
		t.Fatalf("register password validation: %v", err)
	}

	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{name: "valid password", password: "Abcdef1!", valid: true},
		{name: "too short", password: "Ab1!", valid: false},
		{name: "missing uppercase", password: "abcdef1!", valid: false},
		{name: "missing lowercase", password: "ABCDEF1!", valid: false},
		{name: "missing digit", password: "Abcdefgh!", valid: false},
		{name: "missing special", password: "Abcdef12", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.password, "password")
			if tt.valid && err != nil {
				t.Fatalf("expected valid password, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Fatal("expected invalid password")
			}
		})
	}
}
