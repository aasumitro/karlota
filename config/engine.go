package config

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var AllowOrigins = []string{
	"http://localhost:3000",
	"http://localhost:8000",
	"http://192.168.89.82:8000",
}

var allowMethods = []string{"GET, POST, PATCH, DELETE"}

var allowHeaders = []string{
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"Authorization",
	"Cache-Control",
	"Origin",
	"Cookie",
	"X-REFRESH-TOKEN",
}

var exposeHeaders = []string{"Content-Length"}

var allowOriginFunc = func(origin string) bool {
	return origin == "http://localhost:3000"
}

func GinEngine() Option {
	return func(cfg *Config) {
		log.Printf("Trying to init engine (GIN %s) . . . .", gin.Version)
		// set gin mode
		gin.SetMode(func() string {
			if cfg.ServerDebug {
				return gin.DebugMode
			}
			return gin.ReleaseMode
		}())
		// set global variables
		cfg.Infra.GinEngine = gin.Default()
		// set cors middleware
		cfg.Infra.GinEngine.Use(cors.New(cors.Config{
			AllowOrigins:     AllowOrigins,
			AllowMethods:     allowMethods,
			AllowHeaders:     allowHeaders,
			ExposeHeaders:    exposeHeaders,
			AllowCredentials: true,
			AllowOriginFunc:  allowOriginFunc,
			MaxAge:           12 * time.Hour,
		}))
	}
}
