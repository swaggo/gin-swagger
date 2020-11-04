# gin-swagger

gin middleware to automatically generate RESTful API documentation with Swagger 2.0.

[![Travis branch](https://img.shields.io/travis/swaggo/gin-swagger/master.svg)](https://travis-ci.org/swaggo/gin-swagger)
[![Codecov branch](https://img.shields.io/codecov/c/github/swaggo/gin-swagger/master.svg)](https://codecov.io/gh/swaggo/gin-swagger)
[![Go Report Card](https://goreportcard.com/badge/github.com/swaggo/gin-swagger)](https://goreportcard.com/report/github.com/swaggo/gin-swagger)
[![GoDoc](https://godoc.org/github.com/swaggo/gin-swagger?status.svg)](https://godoc.org/github.com/swaggo/gin-swagger)


## Usage

### Start using it
1. Add comments to your API source code, [See Declarative Comments Format](https://swaggo.github.io/swaggo.io/declarative_comments_format/).
2. Download [Swag](https://github.com/swaggo/swag) for Go by using:
```sh
$ go get -u github.com/swaggo/swag/cmd/swag
```

3. Run the [Swag](https://github.com/swaggo/swag) in your Go project root folder which contains `main.go` file, [Swag](https://github.com/swaggo/swag) will parse comments and generate required files(`docs` folder and `docs/doc.go`).
```sh
$ swag init
```
4. Download [gin-swagger](https://github.com/swaggo/gin-swagger) by using:
```sh
$ go get -u github.com/swaggo/gin-swagger
$ go get -u github.com/swaggo/files
```
And import following in your code:

```go
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files" // swagger embed files

```

### Canonical example:
Now assume you have implemented a simple api as following:
```go
//A get function which returns a Helloworld for you 
func Helloworld(g *gin.Context)  {
	g.JSON(http.StatusOK,"helloworld")
}

```
So how to use gin-swagger on api above? Just follow the folling guide:)  

1. Add Comments for apis and main function with gin-swagger rules like following:  
```go
// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context)  {
	g.JSON(http.StatusOK,"helloworld")
}


//Comments used to generate swagger docs are igonred
//We foucs on how to setup gin router in main
func main() {
	r := gin.Default()

    //Using this 
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Controller.Helloworld)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
```

2. Use `swag init` to generate a docs
3. import the docs like this:
Assume your project name is `MyGoSwagger`,don't forget the blank identifier!
```go
import (
    _ "MyGoSwagger/docs"
)
```
4. run your code,and browse to http://localhost:8080/swagger/index.html,you can see Swagger 2.0 Api documents.