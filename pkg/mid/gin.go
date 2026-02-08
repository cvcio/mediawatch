package mid

import (
	"github.com/gin-gonic/gin"
)

// EnableCORS : Enable CORS Middleware in DEV Environment
func EnableCORS(environment string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origins := "*"
		if environment == "production" {
			origins = "https://*.cvcio.com, https://*.mediawatch.io"
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origins)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", " X-Requested-With, Content-Type, Origin, Accept, Client-Security-Token, Accept-Encoding, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
