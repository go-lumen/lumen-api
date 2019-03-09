package tests

import (
	"os"
	"testing"

	"github.com/adrien3d/base-api/models"
	"github.com/adrien3d/base-api/server"
)

var api *server.API
var user *models.User
var authToken string

func TestMain(m *testing.M) {
	api = SetupApi()
	user, authToken = CreateUserAndGenerateToken()
	retCode := m.Run()
	api.Database.Session.Close()
	os.Exit(retCode)
}
