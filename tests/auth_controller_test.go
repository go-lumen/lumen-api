package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	parameters := []byte(`
	{
		"email":"adrien@plugblocks.com",
		"password":"adminpwd"
	}`)

	resp := SendRequest(parameters, "POST", "/v1/auth/")
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestLogOut(t *testing.T) {
	defer ResetDatabase()

	resp := SendRequestWithToken(nil, "GET", "/v1/auth/logout", authToken)
	assert.Equal(t, http.StatusOK, resp.Code)

	resp = SendRequestWithToken(nil, "GET", "/v1/users/"+user.ID, authToken)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
