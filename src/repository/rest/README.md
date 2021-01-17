## REST

### users_repository.go

rest client library
- github.com/cmd-ctrl-q/golang-restclient/rest


internal/private url accessed only through internal vpn
no one in the external world can hit this url
`response := usersRestClient.Post("/users/login", request)`


### users_repository_test.go

You want to be able to mock every possible response from the api call. 
``` Go
response := usersRestClient.Post("/users/login", request)
```

before running the oauth application, we had all tests run successfully.

Test code 
``` Go
package rest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/cmd-ctrl-q/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

// as the entry point for an app is the main function,
// the entry point for every test case, is the TestMain() function
// on the package 'rest'
func TestMain(m *testing.M) {
	fmt.Println("about to start test cases...")
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromAPI(t *testing.T) {
	// remove any mocks from previous test cases.
	rest.FlushMockups()
	// configure mocks
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@email.com","password":"password"}`,
		RespHTTPCode: -1,   // we want invalid http code
		RespBody:     `{}`, // not necessary
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@email.com", "password")
	// fmt.Println(user)
	// fmt.Println(err)

	assert.Nil(t, user)   // expect user to be nil
	assert.NotNil(t, err) // err not nil bc were expecting an err
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restclient response when tyring to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	// remove any mocks from previous test cases.
	rest.FlushMockups()
	// configure mocks
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@email.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound, // we want invalid http code
		RespBody:     `{"message": "invalid login credentials", "status": "404", "error": "not_found"}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@email.com", "password")
	// fmt.Println(user)
	// fmt.Println(err)

	assert.Nil(t, user)   // expect user to be nil
	assert.NotNil(t, err) // err not nil bc were expecting an err
	// we expect a status internal server error because the above
	// request did not get unmarshalled.
	// as there was an error unmarshalling it.
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	// remove any mocks from previous test cases.
	rest.FlushMockups()
	// configure mocks
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@email.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound, // we want invalid http code
		RespBody:     `{"message": "invalid login credentials", "status": 404, "error": "not_found"}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@email.com", "password")
	// fmt.Println(user)
	// fmt.Println(err)

	assert.Nil(t, user)   // expect user to be nil
	assert.NotNil(t, err) // err not nil bc were expecting an err
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	// the reason we expect "invalid login creds..." is because the
	// request was successfully unmarshalled.
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

// error when unmarshalling a errors.RestErr
func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	// remove any mocks from previous test cases.
	rest.FlushMockups()
	// configure mocks
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@email.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "6","first_name": "wolf","last_name": "wolf","email": "wolfwwolf@email.com"}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@email.com", "password")
	// fmt.Println(user)
	// fmt.Println(err)

	assert.Nil(t, user)   // expect user to be nil
	assert.NotNil(t, err) // err not nil bc were expecting an err
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	// the reason we expect "invalid login creds..." is because the
	// error request was successfully unmarshalled.
	assert.EqualValues(t, "error when trying to unmarshal users login reponse", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	// remove any mocks from previous test cases.
	rest.FlushMockups()
	// configure mocks
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@email.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 6,"first_name": "wolf","last_name": "wolf","email": "wolfwwolf@email.com"}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("wolfwwolf@email.com", "woof")
	// fmt.Println(user)
	// fmt.Println(err)

	assert.Nil(t, err)     // expect nerr to be nil
	assert.NotNil(t, user) // user not nil bc were expecting a user
	assert.EqualValues(t, 6, user.ID)
	assert.EqualValues(t, "wolf", user.FirstName)
	assert.EqualValues(t, "wolf", user.LastName)
	assert.EqualValues(t, "wolfwwolf@email.com", user.Email)
}
```

