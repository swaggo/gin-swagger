package ginSwagger

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/swaggo/gin-swagger/example/basic/docs"
)

func TestWrapHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/swagger/*any", WrapHandler(swaggerFiles.Handler))

	w1 := performRequest("GET", "/swagger/", router)
	assert.Equal(t, 200, w1.Code)
}

func TestCustomWrapHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/*any", CustomWrapHandler(&Config{SwaggerBase: "/"}, swaggerFiles.Handler))

	w1 := performRequest("GET", "/", router)
	assert.Equal(t, 200, w1.Code)

	w2 := performRequest("GET", "/doc.json", router)
	assert.Equal(t, 200, w2.Code)

	w3 := performRequest("GET", "/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/oauth2-redirect.html", router)
	assert.Equal(t, 200, w4.Code)

	w5 := performRequest("GET", "/notfound", router)
	assert.Equal(t, 404, w5.Code)
}

func TestDisablingWrapHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	disablingKey := "SWAGGER_DISABLE"

	router.GET("/simple/*any", DisablingWrapHandler(swaggerFiles.Handler, disablingKey, SwaggerBase("/simple/")))

	w1 := performRequest("GET", "/simple/", router)
	assert.Equal(t, 200, w1.Code)

	w2 := performRequest("GET", "/simple/doc.json", router)
	assert.Equal(t, 200, w2.Code)

	w3 := performRequest("GET", "/simple/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/simple/oauth2-redirect.html", router)
	assert.Equal(t, 200, w4.Code)

	w5 := performRequest("GET", "/simple/notfound", router)
	assert.Equal(t, 404, w5.Code)

	os.Setenv(disablingKey, "true")

	router.GET("/disabling/*any", DisablingWrapHandler(swaggerFiles.Handler, disablingKey, SwaggerBase("/disabling/")))

	w11 := performRequest("GET", "/disabling/", router)
	assert.Equal(t, 404, w11.Code)

	w22 := performRequest("GET", "/disabling/doc.json", router)
	assert.Equal(t, 404, w22.Code)

	w33 := performRequest("GET", "/disabling/favicon-16x16.png", router)
	assert.Equal(t, 404, w33.Code)

	w44 := performRequest("GET", "/disabling/oauth2-redirect.html", router)
	assert.Equal(t, 404, w44.Code)

	w55 := performRequest("GET", "/disabling/notfound", router)
	assert.Equal(t, 404, w55.Code)
}

func TestDisablingCustomWrapHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	disablingKey := "SWAGGER_DISABLE2"

	router.GET("/simple/*any", DisablingCustomWrapHandler(&Config{SwaggerBase: "/simple/"}, swaggerFiles.Handler, disablingKey))

	w1 := performRequest("GET", "/simple/", router)
	assert.Equal(t, 200, w1.Code)

	os.Setenv(disablingKey, "true")

	router.GET("/disabling/*any", DisablingCustomWrapHandler(&Config{SwaggerBase: "/disabling/"}, swaggerFiles.Handler, disablingKey))

	w11 := performRequest("GET", "/disabling/", router)
	assert.Equal(t, 404, w11.Code)
}

func TestWithGzipMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(gzip.Gzip(gzip.BestSpeed))

	router.GET("/*any", WrapHandler(swaggerFiles.Handler, SwaggerBase("/")))

	w1 := performRequest("GET", "/", router)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, w1.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w2 := performRequest("GET", "/swagger-ui.css", router)
	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, w2.Header()["Content-Type"][0], "text/css; charset=utf-8")

	w3 := performRequest("GET", "/swagger-ui-bundle.js", router)
	assert.Equal(t, 200, w3.Code)
	assert.Equal(t, w3.Header()["Content-Type"][0], "application/javascript")

	w4 := performRequest("GET", "/doc.json", router)
	assert.Equal(t, 200, w4.Code)
	assert.Equal(t, w4.Header()["Content-Type"][0], "application/json")
}

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
