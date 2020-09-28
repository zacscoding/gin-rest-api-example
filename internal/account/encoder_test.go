package account

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncodePassword(t *testing.T) {
	password := "password1234"
	// when
	encoded, err := EncodePassword(password)

	// then
	assert.NoError(t, err)
	err = MatchesPassword(encoded, password)
	assert.NoError(t, err)
}

func TestMatchesPassword(t *testing.T) {
	password := "password1234"
	encoded, err := EncodePassword(password)
	assert.NoError(t, err)

	// when then
	err = MatchesPassword(encoded, password)
	assert.NoError(t, err)

	// when then
	err = MatchesPassword(encoded, password+"append")
	assert.Error(t, err)
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
}
