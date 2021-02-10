package rest

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/domain/users"
	"github.com/cmd-ctrl-q/bookstore_utils-go/rest_errors"
	"github.com/cmd-ctrl-q/golang-restclient/rest"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8081", // oauth api is consuming users api which is on port 8081
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, rest_errors.RestErr)
}

type usersRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, rest_errors.RestErr) {

	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	// bytes, _ := json.Marshal(request)
	// fmt.Println(string(bytes))

	response := usersRestClient.Post("/users/login", request) // internal/private url accessed only through internal vpn

	// rest client timeout
	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when tyring to login user", errors.New("restclient error"))
	}

	// any other error.
	// invalid error whose struct signature doesnt match our restErr fields
	// if status code > 299 respond with a rest error to caller
	if response.StatusCode > 299 {
		// purpose: use response.Bytes() (ie response json from the api) to create a new RestErr
		apiErr, err := rest_errors.NewRestErrorFromBytes(response.Bytes())
		// err := json.Unmarshal(response.Bytes(), restErr)

		// if there err != nil, then someone change the signature of the error response in the users api
		// and now its not responding a rest error
		if err != nil {
			// if the json from the response is valid then itll return just an internal server error
			// ie idk what to do with the request
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, apiErr
	}
	// success bc status code < 299
	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal users login reponse", errors.New("json parsing error"))
	}
	return &user, nil
}
