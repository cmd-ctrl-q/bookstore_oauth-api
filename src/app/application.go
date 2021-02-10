package app

import (
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/http"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/repository/db"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/repository/rest"
	"github.com/cmd-ctrl-q/bookstore_oauth-api/src/services/access_token"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// StartApplication starts the application
func StartApplication() {
	// session := cassandra.GetSession()
	// session.Close()

	atHandler := http.NewAccessTokenHandler(access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetByID)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
