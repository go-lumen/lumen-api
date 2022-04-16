package tests

import (
	"net/http"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	//Missing field
	parameters := []byte(`
	{
		"password":"adminpwd",
		"first_name":"Adrien",
		"last_name": "Chapelet"
	}`)
	resp := SendRequest(parameters, "POST", "/v1/users/")
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	//Everything is fine
	parameters = []byte(`
	{
		"email":"adrien@plugblocks.com",
		"password":"adminpwd",
		"first_name":"Adrien",
		"last_name": "Chapelet"
	}`)
	resp = SendRequest(parameters, "POST", "/v1/users/")
	assert.Equal(t, http.StatusCreated, resp.Code)

	// User already exists
	resp = SendRequest(parameters, "POST", "/v1/users/")
	assert.Equal(t, http.StatusConflict, resp.Code)

	// Duplicate email
	parameters = []byte(`
	{
		"email":"aDrIeN@plugblocks.com",
		"password":"adminpwd",
		"first_name":"Adrien",
		"last_name": "Chapelet"
	}`)
	resp = SendRequest(parameters, "POST", "/v1/users/")
	assert.Equal(t, http.StatusConflict, resp.Code)

	// Test activation
	/*user := models.User{}
	err := api.MongoDatabase.C(models.UsersCollection).Find(bson.M{"email": "adrien@plugblocks.com"}).One(&user)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, user.Active, false)
	resp = SendRequest(nil, "GET", "/v1/users/"+user.ID+"/activate/"+user.ActivationKey)

	//Update user information
	err = api.MongoDatabase.C(models.UsersCollection).Find(bson.M{"email": "adrien@plugblocks.com"}).One(&user)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, user.Active, true)*/

	//Activation key isn't right
	resp = SendRequest(nil, "GET", "/v1/users/"+user.ID+"/activate/fakeKey")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	//Unknown user
	resp = SendRequest(nil, "GET", "/v1/users/"+bson.NewObjectId().Hex()+"/activate/fakeKey")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
