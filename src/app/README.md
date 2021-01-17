


### 

`/oauth/access_token/:access_token_id`
you can make this get request to `/oauth/access_token/:access_token_id` from any programming language you want.
thus, you can create a shared library in the go language to use in every project using go.


NOTICE: 
- all external requests must first go through the OAuth API before they can use other APIs

OAuth API sends to Uses API
- request json
``` JSON
{
	"email": "emailaddr@email.com",
	"password": "123abc"
}
```
OAuth API sends to Uses API
- response json
``` JSON
{
	"grant_type": "password",
	"username": "emailaddr@email.com",
	"password": "123abc"
}
```
grant_type: e.g. client_credentials: client_id, client_secret 
 - client exchanges client_id and client_secret for the access token. 
``` JSON
{
	"grant_type": "client_credentials",
	"client_id": "id-123",
	"client_secret": "secret-123"
}
```

### Old app code 

``` Go
package app

import (
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/clients/cassandra"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/domain/access_token"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/http"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/repository/db"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	session := cassandra.GetSession()
	// if dbErr != nil {
	// 	panic(dbErr)
	// }
	session.Close()

	atHandler := http.NewHandler(access_token.NewService(db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetByID)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
```