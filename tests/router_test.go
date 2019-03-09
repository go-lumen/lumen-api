package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/", nil)
	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}
