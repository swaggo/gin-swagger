# gin-swagger

gin middleware to automatically generate RESTful API documentation with Swagger 2.0.

[![Travis branch](https://img.shields.io/travis/swag-gonic/gin-swagger/master.svg)](https://travis-ci.org/swag-gonic/gin-swagger)
[![Codecov branch](https://img.shields.io/codecov/c/github/swag-gonic/gin-swagger/master.svg)](https://codecov.io/gh/swag-gonic/gin-swagger)
[![Go Report Card](https://goreportcard.com/badge/github.com/swag-gonic/gin-swagger)](https://goreportcard.com/report/github.com/swag-gonic/gin-swagger)


## Usage

### Start using it
1. Add comments to your API source code, [see Declarative Comments Format](#declarative-comments-format)
2. Download Swag for Go by using:
```sh
$ go get -u github.com/swag-gonic/swag
```
3. Run the Swag in your Go project root folder which contains `main.go` file, Swag will parse your comments and generate required files(`docs` folder and `docs/doc.go`)
```sh
$ swag init
```
3.
```

```sh
$ go get github.com/swag-gonic/gin-swagger
```

Import it in your code:

```go
import "github.com/gin-gonic/gin"
import "github.com/swag-gonic/gin-swagger"
import "github.com/swag-gonic/gin-swagger/swaggerFiles"

```

### Canonical example:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swag-gonic/gin-swagger"
	"github.com/swag-gonic/gin-swagger/swaggerFiles"

	_ "github.com/swag-gonic/gin-swagger/example/docs"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
```