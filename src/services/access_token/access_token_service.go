package access_token

import (
	"strings"

	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/domain/access_token"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/repository/db"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/repository/rest"
	"github.com/cmd-ctrl-q/bookstore_utils-go/rest_errors"
)

type Repository interface {
	GetByID(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessToken) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type Service interface {
	GetByID(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetByID(accessTokenID string) (*access_token.AccessToken, rest_errors.RestErr) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}

	// valid accessTokenID
	accessToken, err := s.dbRepo.GetByID(accessTokenID)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

// Create takes in a newly created access token as a paramater before creating it
func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr) {

	// validate the request
	if err := request.Validate(); err != nil {
		return nil, err
	}

	//TODO: support both grant types: client_Credentials and password
	// currently only supporting password but could support client_credentials

	// ----- create access tokens based on grant_type password
	// authenticate the user against the users api
	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// generate a new access token
	at := access_token.GetNewAccessToken(user.ID)
	at.Generate()
	// ------ end

	// save the new access token to cassandra
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
