package rest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/domain/users"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/utils/errors"
	"github.com/cmd-ctrl-q/golang-restclient/rest"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8081", // oauth api is consuming users api which is on port 8081
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {

	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	// bytes, _ := json.Marshal(request)
	// fmt.Println(string(bytes))

	response := usersRestClient.Post("/users/login", request) // internal/private url accessed only through internal vpn

	// rest client timeout
	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServerError("invalid restclient response when tyring to login user")
	}

	// any other error.
	// invalid error whose struct signature doesnt match our restErr fields
	if response.StatusCode > 299 {
		fmt.Println(response.String())
		var restErr errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		// if there err != nil, then someone change the signature of the error response in the users api
		// and now its not responding a rest error
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}
	// success bc status code < 299
	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal users login reponse")
	}
	return &user, nil
}
