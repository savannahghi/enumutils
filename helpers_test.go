package go_utils_test

import (
	"testing"

	base "github.com/savannahghi/go_utils"
)

func TestGenerateRandomWithNDigits(t *testing.T) {
	result, err := base.GenerateRandomWithNDigits(5)
	if result == "" {
		t.Errorf("can't generate random with n digits")
		return
	}
	if err != nil {
		t.Errorf("can't generate random with n digits: %v", err)
		return
	}
}

func TestGenerateRandomEmail(t *testing.T) {
	email := base.GenerateRandomEmail()
	if email == "" {
		t.Errorf("can't generate a unique email")
	}
}
