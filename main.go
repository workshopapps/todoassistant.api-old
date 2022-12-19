package main

import (
	"test-va/cmd"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// @title           Ticked
	// @version         1.0
	// @description     A web application that helps you to keep your day planned
	// @termsOfService  http://swagger.io/terms/

	// @contact.name   API Support
	// @contact.url    http://www.swagger.io/support
	// @contact.email  support@swagger.io

	// @license.name  Apache 2.0
	// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

	// @host      api.ticked.hng.tech:2022
	// @BasePath  /api/v1

	// @securityDefinitions.apiKey  ApiKeyAuth
	// @in header
	// @name Authorisation
	cmd.Setup()
}
