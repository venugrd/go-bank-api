package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	pass := RandomString(6)

	hash, err := HashPassword(pass)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = IsValidPassword(hash, pass)
	require.NoError(t, err)

	err = IsValidPassword(hash, RandomString(6))
	require.Error(t, err)
	fmt.Println(err)
}
