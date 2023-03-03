// CORS跨域处理中间件
package middleware

import (
	"net/http"
	"net/url"
	"os"

	"github.com/perpower/goframe/funcs/convert"
	"github.com/perpower/goframe/funcs/normal"
	"github.com/perpower/goframe/funcs/pos"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// 服务端允许跨域请求选项
type CorsOptions struct {
	AllowDomain      []string `yaml:"allowDomain" json:"allowDomain"`           // Used for allowing requests from custom domains
	AllowOrigin      string   `yaml:"allowOrigin" json:"allowOrigin"`           // Access-Control-Allow-Origin
	AllowCredentials string   `yaml:"allowCredentials" json:"allowCredentials"` // Access-Control-Allow-Credentials
	ExposeHeaders    string   `yaml:"exposeHeaders" json:"exposeHeaders"`       // Access-Control-Expose-Headers
	MaxAge           int      `yaml:"maxAge" json:"maxAge"`                     // Access-Control-Max-Age
	AllowMethods     string   `yaml:"allowMethod" json:"allowMethod"`           // Access-Control-Allow-Methods
	AllowHeaders     string   `yaml:"allowHeaders" json:"allowHeaders"`         // Access-Control-Allow-Headers
}

var (
	// defaultAllowHeaders is the default allowed headers for CORS.
	// It defined another map for better header key searching performance.
	defaultAllowHeaders    = "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With"
	defaultAllowHeadersMap = make(map[string]struct{})
	supportedHttpMethods   = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
)

func init() {
	array := normal.SplitAndTrim(defaultAllowHeaders, ",")
	for _, header := range array {
		defaultAllowHeadersMap[header] = struct{}{}
	}
}

func CorsHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := CorsOptions{}
		if f, err := os.ReadFile("./configs/cors.yml"); err != nil {
			conf := DefaultCorsOptions(c)
			SetCors(c, conf)
		} else {
			err = yaml.Unmarshal(f, &conf) //解析yaml文件配置内容
			if err != nil {
				panic(err)
			}
			SetCors(c, conf)
		}

		c.Next()
	}
}

// DefaultCorsOptions returns the default CORS options,
// which allows any cross-domain request.
func DefaultCorsOptions(c *gin.Context) CorsOptions {
	options := CorsOptions{
		AllowOrigin:      "*",
		AllowMethods:     supportedHttpMethods,
		AllowCredentials: "true",
		AllowHeaders:     defaultAllowHeaders,
		MaxAge:           3628800,
	}
	// Allow all client's custom headers in default.
	if headers := c.Request.Header.Get("Access-Control-Request-Headers"); headers != "" {
		array := normal.SplitAndTrim(headers, ",")
		for _, header := range array {
			if _, ok := defaultAllowHeadersMap[header]; !ok {
				options.AllowHeaders += "," + header
			}
		}
	}
	// Allow all anywhere origin in default.
	if origin := c.Request.Header.Get("Origin"); origin != "" {
		options.AllowOrigin = origin
	} else if referer := c.Request.Referer(); referer != "" {
		if p := pos.PosR(referer, "/", 6); p != -1 {
			options.AllowOrigin = referer[:p]
		} else {
			options.AllowOrigin = referer
		}
	}
	return options
}

// CORS sets custom CORS options.
// See https://www.w3.org/TR/cors/ .
func SetCors(c *gin.Context, options CorsOptions) {
	if CorsAllowedOrigin(c, options) {
		c.Header("Access-Control-Allow-Origin", options.AllowOrigin)
	} else {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if options.AllowCredentials != "" {
		c.Header("Access-Control-Allow-Credentials", options.AllowCredentials)
	}
	if options.ExposeHeaders != "" {
		c.Header("Access-Control-Expose-Headers", options.ExposeHeaders)
	}
	if options.MaxAge != 0 {
		c.Header("Access-Control-Max-Age", convert.String(options.MaxAge))
	}
	if options.AllowMethods != "" {
		c.Header("Access-Control-Allow-Methods", options.AllowMethods)
	}
	if options.AllowHeaders != "" {
		c.Header("Access-Control-Allow-Headers", options.AllowHeaders)
	}
	// No continue service handling if it's OPTIONS request.
	// Note that there's special checks in previous router searching,
	// so if it goes to here it means there's already serving handler exist.
	if normal.Equal(c.Request.Method, "OPTIONS") {
		c.AbortWithStatus(http.StatusNoContent)
		// No continue serving.
		return
	}
}

// CORSAllowedOrigin CORSAllowed checks whether the current request origin is allowed cross-domain.
func CorsAllowedOrigin(c *gin.Context, options CorsOptions) bool {
	if options.AllowDomain == nil || len(options.AllowDomain) == 0 || normal.InArray("*", options.AllowDomain) {
		return true
	}
	origin := c.Request.Header.Get("Origin")
	if origin == "" {
		return true
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if normal.InArray(parsed.Host, options.AllowDomain) {
		return true
	}

	return false
}

// CORSDefault sets CORS with default CORS options,
// which allows any cross-domain request.
func CorsDefault(c *gin.Context) {
	SetCors(c, DefaultCorsOptions(c))
}
