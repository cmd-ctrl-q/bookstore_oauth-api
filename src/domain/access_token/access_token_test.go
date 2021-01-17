package access_token

import (
	"testing"
	"time"

	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/domain/users"
	"github.com/stretchr/testify/assert"
)

func TestAccessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime)
}

func TestGetNewAccessToken(t *testing.T) {
	user := users.User{ID: 123}
	at := GetNewAccessToken(user.ID)

	assert.False(t, at.IsExpired(), "new access token should not be expired")
	assert.EqualValues(t, "", at.AccessToken, "new access token should not have defined access token id")
	assert.True(t, at.UserID == 0, "new access token should not have an associated user id")
}

func TestAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}

	assert.True(t, at.IsExpired(), "empty access token should be expired by default")
	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "access expiring three hours from now should NOT have expired")
}
