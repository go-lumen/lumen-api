package tests

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"

	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/server"
	"github.com/spf13/viper"
)

func SendRequest(parameters []byte, method string, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(parameters))
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	return resp
}

func SendRequestWithToken(parameters []byte, method string, url string, authToken string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(parameters))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	return resp
}

func CreateUserAndGenerateToken() (*models.User, string) {
	/*users := api.SetupMongoDatabase.C(models.UsersCollection)

	user := models.User{
		ID:        bson.NewObjectId().Hex(),
		Email:     "adrien@plugblocks.com",
		FirstName: "Adrien",
		LastName:  "Chapelet",
		Password:  "adminpwd",
		Active:    true,
		Admin:     true,
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	users.Insert(user)

	privateKeyFile, _ := ioutil.ReadFile(api.Config.GetString("rsa_private"))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)

	token := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodRS256.Alg()))

	claims := make(jwt.MapClaims)
	// TODO: ADD EXPIRATION
	//claims["exp"] = time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = user.ID

	token.Claims = claims

	tokenString, _ := token.SignedString(privateKey)

	return &user, tokenString*/
	return nil, ""
}

func ResetDatabase() {
	//api.MongoDatabase.DropDatabase()
	user, authToken = CreateUserAndGenerateToken()
}

// SetupAPI allows to setup the api
func SetupAPI() *server.API {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	err := api.SetupViper()
	if err != nil {
		panic(err)
	}

	_, err = api.SetupMongoDatabase()
	if err != nil {
		panic(err)
	}

	/*api.MongoDatabase.DropDatabase()

	err = api.SetupIndexes()
	if err != nil {
		panic(err)
	}

	api.EmailSender = &services.FakeEmailSender{}*/
	api.SetupRouter()

	return api
}
