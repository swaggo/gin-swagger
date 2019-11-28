package main

import (
	ginSwagger "github.com/canecat/gin-swagger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"

	_ "github.com/canecat/gin-swagger/example/basic/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	r := gin.New()

	swaggerBase := ginSwagger.SwaggerBase("swagger/")
	specFileName := ginSwagger.SpecFileName("doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerBase, specFileName))

	r.Run()
}
