package tests

import (
	"os"
	"testing"

	"github.com/adrien3d/lumen-api/models"
	"github.com/adrien3d/lumen-api/server"
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
