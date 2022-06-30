package ginSwagger

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"golang.org/x/net/webdav"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

type swaggerConfig struct {
	URL                      string
	DocExpansion             string
	Title                    string
	Oauth2RedirectURL        template.JS
	DefaultModelsExpandDepth int
	DeepLinking              bool
	PersistAuthorization     bool
	Oauth2DefaultClientID    string
}

// Config stores ginSwagger configuration variables.
type Config struct {
	// The url pointing to API definition (normally swagger.json or swagger.yaml). Default is `doc.json`.
	URL                      string
	DocExpansion             string
	InstanceName             string
	Title                    string
	DefaultModelsExpandDepth int
	DeepLinking              bool
	PersistAuthorization     bool
	Oauth2DefaultClientID    string
}

func (config Config) toSwaggerConfig() swaggerConfig {
	return swaggerConfig{
		URL:                      config.URL,
		DeepLinking:              config.DeepLinking,
		DocExpansion:             config.DocExpansion,
		DefaultModelsExpandDepth: config.DefaultModelsExpandDepth,
		Oauth2RedirectURL: "`${window.location.protocol}//${window.location.host}$" +
			"{window.location.pathname.split('/').slice(0, window.location.pathname.split('/').length - 1).join('/')}" +
			"/oauth2-redirect.html`",
		Title:                 config.Title,
		PersistAuthorization:  config.PersistAuthorization,
		Oauth2DefaultClientID: config.Oauth2DefaultClientID,
	}
}

// URL presents the url pointing to API definition (normally swagger.json or swagger.yaml).
func URL(url string) func(*Config) {
	return func(c *Config) {
		c.URL = url
	}
}

// DocExpansion list, full, none.
func DocExpansion(docExpansion string) func(*Config) {
	return func(c *Config) {
		c.DocExpansion = docExpansion
	}
}

// DeepLinking set the swagger deep linking configuration.
func DeepLinking(deepLinking bool) func(*Config) {
	return func(c *Config) {
		c.DeepLinking = deepLinking
	}
}

// DefaultModelsExpandDepth set the default expansion depth for models
// (set to -1 completely hide the models).
func DefaultModelsExpandDepth(depth int) func(*Config) {
	return func(c *Config) {
		c.DefaultModelsExpandDepth = depth
	}
}

// InstanceName set the instance name that was used to generate the swagger documents
// Defaults to swag.Name ("swagger").
func InstanceName(name string) func(*Config) {
	return func(c *Config) {
		c.InstanceName = name
	}
}

// PersistAuthorization Persist authorization information over browser close/refresh.
// Defaults to false.
func PersistAuthorization(persistAuthorization bool) func(*Config) {
	return func(c *Config) {
		c.PersistAuthorization = persistAuthorization
	}
}

// Oauth2DefaultClientID set the default client ID used for OAuth2
func Oauth2DefaultClientID(oauth2DefaultClientID string) func(*Config) {
	return func(c *Config) {
		c.Oauth2DefaultClientID = oauth2DefaultClientID
	}
}

// WrapHandler wraps `http.Handler` into `gin.HandlerFunc`.
func WrapHandler(handler *webdav.Handler, options ...func(*Config)) gin.HandlerFunc {
	var config = Config{
		URL:                      "doc.json",
		DocExpansion:             "list",
		InstanceName:             swag.Name,
		Title:                    "Swagger UI",
		DefaultModelsExpandDepth: 1,
		DeepLinking:              true,
		PersistAuthorization:     false,
		Oauth2DefaultClientID:    "",
	}

	for _, c := range options {
		c(&config)
	}

	return CustomWrapHandler(&config, handler)
}

// CustomWrapHandler wraps `http.Handler` into `gin.HandlerFunc`.
func CustomWrapHandler(config *Config, handler *webdav.Handler) gin.HandlerFunc {
	var once sync.Once

	if config.InstanceName == "" {
		config.InstanceName = swag.Name
	}

	if config.Title == "" {
		config.Title = "Swagger UI"
	}

	// create a template with name
	index, _ := template.New("swagger_index.html").Parse(swaggerIndexTpl)

	var matcher = regexp.MustCompile(`(.*)(index\.html|doc\.json|favicon-16x16\.png|favicon-32x32\.png|/oauth2-redirect\.html|swagger-ui\.css|swagger-ui\.css\.map|swagger-ui\.js|swagger-ui\.js\.map|swagger-ui-bundle\.js|swagger-ui-bundle\.js\.map|swagger-ui-standalone-preset\.js|swagger-ui-standalone-preset\.js\.map)[?|.]*`)

	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.AbortWithStatus(http.StatusMethodNotAllowed)

			return
		}

		matches := matcher.FindStringSubmatch(ctx.Request.RequestURI)

		if len(matches) != 3 {
			ctx.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))

			return
		}

		path := matches[2]
		once.Do(func() {
			handler.Prefix = matches[1]
		})

		switch filepath.Ext(path) {
		case ".html":
			ctx.Header("Content-Type", "text/html; charset=utf-8")
		case ".css":
			ctx.Header("Content-Type", "text/css; charset=utf-8")
		case ".js":
			ctx.Header("Content-Type", "application/javascript")
		case ".png":
			ctx.Header("Content-Type", "image/png")
		case ".json":
			ctx.Header("Content-Type", "application/json; charset=utf-8")
		}

		switch path {
		case "index.html":
			_ = index.Execute(ctx.Writer, config.toSwaggerConfig())
		case "doc.json":
			doc, err := swag.ReadDoc(config.InstanceName)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)

				return
			}

			ctx.String(http.StatusOK, doc)
		default:
			handler.ServeHTTP(ctx.Writer, ctx.Request)
		}
	}
}

// DisablingWrapHandler turn handler off
// if specified environment variable passed.
func DisablingWrapHandler(handler *webdav.Handler, envName string) gin.HandlerFunc {
	if os.Getenv(envName) != "" {
		return func(c *gin.Context) {
			// Simulate behavior when route unspecified and
			// return 404 HTTP code
			c.String(http.StatusNotFound, "")
		}
	}

	return WrapHandler(handler)
}

// DisablingCustomWrapHandler turn handler off
// if specified environment variable passed.
func DisablingCustomWrapHandler(config *Config, handler *webdav.Handler, envName string) gin.HandlerFunc {
	if os.Getenv(envName) != "" {
		return func(c *gin.Context) {
			// Simulate behavior when route unspecified and
			// return 404 HTTP code
			c.String(http.StatusNotFound, "")
		}
	}

	return CustomWrapHandler(config, handler)
}

const swaggerIndexTpl = `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>{{.Title}}</title>
  <style>@font-face{font-family:'Open Sans';font-style:normal;font-weight:400;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSKmu0SC55K5gw.woff2) format('woff2');unicode-range:U+0460-052F,U+1C80-1C88,U+20B4,U+2DE0-2DFF,U+A640-A69F,U+FE2E-FE2F;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:400;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSumu0SC55K5gw.woff2) format('woff2');unicode-range:U+0301,U+0400-045F,U+0490-0491,U+04B0-04B1,U+2116;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:400;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSOmu0SC55K5gw.woff2) format('woff2');unicode-range:U+1F00-1FFF;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:400;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSymu0SC55K5gw.woff2) format('woff2');unicode-range:U+0370-03FF;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:400;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTS2mu0SC55K5gw.woff2) format('woff2');unicode-range:U+0590-05FF,U+200C-2010,U+20AA,U+25CC,U+FB1D-FB4F;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:400;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSCmu0SC55K5gw.woff2) format('woff2');unicode-range:U+0102-0103,U+0110-0111,U+0128-0129,U+0168-0169,U+01A0-01A1,U+01AF-01B0,U+1EA0-1EF9,U+20AB;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:400;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSGmu0SC55K5gw.woff2) format('woff2');unicode-range:U+0100-024F,U+0259,U+1E00-1EFF,U+2020,U+20A0-20AB,U+20AD-20CF,U+2113,U+2C60-2C7F,U+A720-A7FF;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:400;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTS-mu0SC55I.woff2) format('woff2');unicode-range:U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+2000-206F,U+2074,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:700;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSKmu0SC55K5gw.woff2) format('woff2');unicode-range:U+0460-052F,U+1C80-1C88,U+20B4,U+2DE0-2DFF,U+A640-A69F,U+FE2E-FE2F;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:700;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSumu0SC55K5gw.woff2) format('woff2');unicode-range:U+0301,U+0400-045F,U+0490-0491,U+04B0-04B1,U+2116;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:700;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSOmu0SC55K5gw.woff2) format('woff2');unicode-range:U+1F00-1FFF;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:700;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSymu0SC55K5gw.woff2) format('woff2');unicode-range:U+0370-03FF;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:700;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTS2mu0SC55K5gw.woff2) format('woff2');unicode-range:U+0590-05FF,U+200C-2010,U+20AA,U+25CC,U+FB1D-FB4F;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:700;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSCmu0SC55K5gw.woff2) format('woff2');unicode-range:U+0102-0103,U+0110-0111,U+0128-0129,U+0168-0169,U+01A0-01A1,U+01AF-01B0,U+1EA0-1EF9,U+20AB;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:700;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTSGmu0SC55K5gw.woff2) format('woff2');unicode-range:U+0100-024F,U+0259,U+1E00-1EFF,U+2020,U+20A0-20AB,U+20AD-20CF,U+2113,U+2C60-2C7F,U+A720-A7FF;}@font-face{font-family:'Open Sans';font-style:normal;font-weight:700;font-stretch:100%;src:url(https://fonts.gstatic.com/s/opensans/v29/memvYaGs126MiZpBA-UvWbX2vVnXBbObj2OVTS-mu0SC55I.woff2) format('woff2');unicode-range:U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+2000-206F,U+2074,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:300;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlMOvWnsUnxlC9.woff2) format('woff2');unicode-range:U+0460-052F,U+1C80-1C88,U+20B4,U+2DE0-2DFF,U+A640-A69F,U+FE2E-FE2F;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:300;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlOevWnsUnxlC9.woff2) format('woff2');unicode-range:U+0301,U+0400-045F,U+0490-0491,U+04B0-04B1,U+2116;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:300;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlMevWnsUnxlC9.woff2) format('woff2');unicode-range:U+1F00-1FFF;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:300;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlPuvWnsUnxlC9.woff2) format('woff2');unicode-range:U+0370-03FF;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:300;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlMuvWnsUnxlC9.woff2) format('woff2');unicode-range:U+0102-0103,U+0110-0111,U+0128-0129,U+0168-0169,U+01A0-01A1,U+01AF-01B0,U+1EA0-1EF9,U+20AB;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:300;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlM-vWnsUnxlC9.woff2) format('woff2');unicode-range:U+0100-024F,U+0259,U+1E00-1EFF,U+2020,U+20A0-20AB,U+20AD-20CF,U+2113,U+2C60-2C7F,U+A720-A7FF;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:300;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlPevWnsUnxg.woff2) format('woff2');unicode-range:U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+2000-206F,U+2074,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlMOvWnsUnxlC9.woff2) format('woff2');unicode-range:U+0460-052F,U+1C80-1C88,U+20B4,U+2DE0-2DFF,U+A640-A69F,U+FE2E-FE2F;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlOevWnsUnxlC9.woff2) format('woff2');unicode-range:U+0301,U+0400-045F,U+0490-0491,U+04B0-04B1,U+2116;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlMevWnsUnxlC9.woff2) format('woff2');unicode-range:U+1F00-1FFF;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlPuvWnsUnxlC9.woff2) format('woff2');unicode-range:U+0370-03FF;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlMuvWnsUnxlC9.woff2) format('woff2');unicode-range:U+0102-0103,U+0110-0111,U+0128-0129,U+0168-0169,U+01A0-01A1,U+01AF-01B0,U+1EA0-1EF9,U+20AB;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlM-vWnsUnxlC9.woff2) format('woff2');unicode-range:U+0100-024F,U+0259,U+1E00-1EFF,U+2020,U+20A0-20AB,U+20AD-20CF,U+2113,U+2C60-2C7F,U+A720-A7FF;}@font-face{font-family:'Source Code Pro';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/sourcecodepro/v21/HI_SiYsKILxRpg3hIP6sJ7fM7PqlPevWnsUnxg.woff2) format('woff2');unicode-range:U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+2000-206F,U+2074,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD;}@font-face{font-family:'Titillium Web';font-style:normal;font-weight:400;src:url(https://fonts.gstatic.com/s/titilliumweb/v15/NaPecZTIAOhVxoMyOr9n_E7fdM3mDaZRbryhsA.woff2) format('woff2');unicode-range:U+0100-024F,U+0259,U+1E00-1EFF,U+2020,U+20A0-20AB,U+20AD-20CF,U+2113,U+2C60-2C7F,U+A720-A7FF;}@font-face{font-family:'Titillium Web';font-style:normal;font-weight:400;src:url(https://fonts.gstatic.com/s/titilliumweb/v15/NaPecZTIAOhVxoMyOr9n_E7fdMPmDaZRbrw.woff2) format('woff2');unicode-range:U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+2000-206F,U+2074,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD;}@font-face{font-family:'Titillium Web';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/titilliumweb/v15/NaPDcZTIAOhVxoMyOr9n_E7ffBzCGIVzY5abuWIGxA.woff2) format('woff2');unicode-range:U+0100-024F,U+0259,U+1E00-1EFF,U+2020,U+20A0-20AB,U+20AD-20CF,U+2113,U+2C60-2C7F,U+A720-A7FF;}@font-face{font-family:'Titillium Web';font-style:normal;font-weight:600;src:url(https://fonts.gstatic.com/s/titilliumweb/v15/NaPDcZTIAOhVxoMyOr9n_E7ffBzCGItzY5abuWI.woff2) format('woff2');unicode-range:U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+2000-206F,U+2074,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD;}@font-face{font-family:'Titillium Web';font-style:normal;font-weight:700;src:url(https://fonts.gstatic.com/s/titilliumweb/v15/NaPDcZTIAOhVxoMyOr9n_E7ffHjDGIVzY5abuWIGxA.woff2) format('woff2');unicode-range:U+0100-024F,U+0259,U+1E00-1EFF,U+2020,U+20A0-20AB,U+20AD-20CF,U+2113,U+2C60-2C7F,U+A720-A7FF;}@font-face{font-family:'Titillium Web';font-style:normal;font-weight:700;src:url(https://fonts.gstatic.com/s/titilliumweb/v15/NaPDcZTIAOhVxoMyOr9n_E7ffHjDGItzY5abuWI.woff2) format('woff2');unicode-range:U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+2000-206F,U+2074,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD;}</style>
  <link rel="stylesheet" type="text/css" href="./swagger-ui.css" >
  <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  <style>
    html
    {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
    }
    *,
    *:before,
    *:after
    {
        box-sizing: inherit;
    }

    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>

<body>

<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="position:absolute;width:0;height:0">
  <defs>
    <symbol viewBox="0 0 20 20" id="unlocked">
          <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V6h2v-.801C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8z"></path>
    </symbol>

    <symbol viewBox="0 0 20 20" id="locked">
      <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8zM12 8H8V5.199C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="close">
      <path d="M14.348 14.849c-.469.469-1.229.469-1.697 0L10 11.819l-2.651 3.029c-.469.469-1.229.469-1.697 0-.469-.469-.469-1.229 0-1.697l2.758-3.15-2.759-3.152c-.469-.469-.469-1.228 0-1.697.469-.469 1.228-.469 1.697 0L10 8.183l2.651-3.031c.469-.469 1.228-.469 1.697 0 .469.469.469 1.229 0 1.697l-2.758 3.152 2.758 3.15c.469.469.469 1.229 0 1.698z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow">
      <path d="M13.25 10L6.109 2.58c-.268-.27-.268-.707 0-.979.268-.27.701-.27.969 0l7.83 7.908c.268.271.268.709 0 .979l-7.83 7.908c-.268.271-.701.27-.969 0-.268-.269-.268-.707 0-.979L13.25 10z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow-down">
      <path d="M17.418 6.109c.272-.268.709-.268.979 0s.271.701 0 .969l-7.908 7.83c-.27.268-.707.268-.979 0l-7.908-7.83c-.27-.268-.27-.701 0-.969.271-.268.709-.268.979 0L10 13.25l7.418-7.141z"/>
    </symbol>


    <symbol viewBox="0 0 24 24" id="jump-to">
      <path d="M19 7v4H5.83l3.58-3.59L8 6l-6 6 6 6 1.41-1.41L5.83 13H21V7z"/>
    </symbol>

    <symbol viewBox="0 0 24 24" id="expand">
      <path d="M10 18h4v-2h-4v2zM3 6v2h18V6H3zm3 7h12v-2H6v2z"/>
    </symbol>

  </defs>
</svg>

<div id="swagger-ui"></div>

<script src="./swagger-ui-bundle.js"> </script>
<script src="./swagger-ui-standalone-preset.js"> </script>
<script>
window.onload = function() {
  // Build a system
  const ui = SwaggerUIBundle({
    url: "{{.URL}}",
    dom_id: '#swagger-ui',
    validatorUrl: null,
    oauth2RedirectUrl: {{.Oauth2RedirectURL}},
    persistAuthorization: {{.PersistAuthorization}},
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
	layout: "StandaloneLayout",
    docExpansion: "{{.DocExpansion}}",
	deepLinking: {{.DeepLinking}},
	defaultModelsExpandDepth: {{.DefaultModelsExpandDepth}}
  })

  const defaultClientId = "{{.Oauth2DefaultClientID}}";
  if (defaultClientId) {
    ui.initOAuth({
      clientId: defaultClientId
    })
  }

  window.ui = ui
}
</script>
</body>

</html>
`
