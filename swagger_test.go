package ginSwagger

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/swaggo/gin-swagger/example/docs"
)

func TestWrapHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/*any", WrapHandler(swaggerFiles.Handler))

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)

	w2 := performRequest("GET", "/doc.json", router)
	assert.Equal(t, 200, w2.Code)

	w3 := performRequest("GET", "/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/notfound", router)
	assert.Equal(t, 404, w4.Code)
}

func TestDisablingWrapHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	disablingKey := "SWAGGER_DISABLE"

	router.GET("/simple/*any", DisablingWrapHandler(swaggerFiles.Handler, disablingKey))

	w1 := performRequest("GET", "/simple/index.html", router)
	assert.Equal(t, 200, w1.Code)

	w2 := performRequest("GET", "/simple/doc.json", router)
	assert.Equal(t, 200, w2.Code)

	w3 := performRequest("GET", "/simple/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/simple/notfound", router)
	assert.Equal(t, 404, w4.Code)

	os.Setenv(disablingKey, "true")

	router.GET("/disabling/*any", DisablingWrapHandler(swaggerFiles.Handler, disablingKey))

	w11 := performRequest("GET", "/disabling/index.html", router)
	assert.Equal(t, 404, w11.Code)

	w22 := performRequest("GET", "/disabling/doc.json", router)
	assert.Equal(t, 404, w22.Code)

	w33 := performRequest("GET", "/disabling/favicon-16x16.png", router)
	assert.Equal(t, 404, w33.Code)

	w44 := performRequest("GET", "/disabling/notfound", router)
	assert.Equal(t, 404, w44.Code)
}

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
