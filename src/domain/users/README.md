



This is the request that the OAuth API sends to the User API 
``` Go
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
```