## OAuth

auth resources
[1](https://darutk.medium.com/oauth-access-token-implementation-30c2e8b90ff0)
[2](https://www.oauth.com/oauth2-servers/access-tokens/access-token-response/)
[3](https://www.oauth.com/oauth2-servers/access-tokens/)



### access_token.go

// Web frontend - Client-ID: 123
// 	maybe you want to limit the access token to some api's available for only web applications.
// Android APP - Client-ID: 234
// 	android may have longer access token expiration time


Time stamp of the current time in UTC with 24 hours added to it.
``` Go
Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
```

``` Go
If now is after the expiration time, then it expired
func (at AccessToken) IsExpired() bool {
	now := time.Now().UTC()
	expirationTime := time.Unix(at.Expires, 0)

    // ie if now is after expiration time, return true
	return now.After(expirationTime)
}
```

### access_token_test.go

If access token = 0, then it already expired
``` Go
func TestAccessTokenIsExpired(t *testing.T) {
    at := AccessToken{}
    // at = 0 
	if !at.IsExpired() {
		t.Error("empty access token should be expired by default")
	}
}
```
``` Go
func TestAccessTokenConstants(t *testing.T) {
    // option 1
	if expirationTime != 24 {
		t.Error("expiration time should be 24 hours")
    }
    // option 2
	assert.EqualValues(t, 24, expirationTime, "expiration time should be 24 hours")
}
```

``` Go
func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken()
	if at.IsExpired() {
		t.Error("new access token should not be expired")
	}
	// assume is expired is false
	assert.False(t, at.IsExpired(), "new access token should not be expired")

	if at.AccessToken != "" {
		t.Error("new access token should not have defined access token id")
	}
	// assume at is empty
	assert.EqualValues(t, "", at.AccessToken, "new access token should not have defined access token id")

	if at.UserID != 0 {
		t.Error("new access token should not have an associated user id")
	}
	assert.True(t, at.UserID == 0, "new access token should not have an associated user id")
	assert.EqualValues(t, 0, at.UserID, "new access token should not have an associated user id")
}

func TestAccessTokenIsExpired(t *testing.T) {
    at := AccessToken{}
    // these 3 are same
	// if IsExpired is false, ie not expired, then it should not have expired
	if !at.IsExpired() {
		t.Error("empty access token should be expired by default")
	}
	assert.False(t, !at.IsExpired(), "empty access token should be expired by default")
	assert.True(t, at.IsExpired(), "empty access token should be expired by default")


    at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
    
    // these 2 are same
	// if true, at expired
	if at.IsExpired() {
		t.Error("access expiring three hours from now should NOT have expired")
	}
	assert.False(t, at.IsExpired(), "access expiring three hours from now should NOT have expired")
}
```

``` Go
func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken()
	// assume isExpired is false
	assert.False(t, at.IsExpired(), "new access token should not be expired")

	// assume at is empty, if not empty, display msg
	assert.EqualValues(t, "", at.AccessToken, "new access token should not have defined access token id")

	// assume userID is 0, ie no user, if not 0 display msg
	assert.True(t, at.UserID == 0, "new access token should not have an associated user id")
}

func TestAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}
	// assume at expired is true, if it did not expire, then display msg
	assert.True(t, at.IsExpired(), "empty access token should be expired by default")

	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	// assume at expired is false, so if it did expire, then display msg
	assert.False(t, at.IsExpired(), "access expiring three hours from now should NOT have expired")
}
```

