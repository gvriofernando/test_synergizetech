package header

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSAllowHeaders() gin.HandlerFunc {
	cf := cors.DefaultConfig()
	cf.AllowAllOrigins = true
	cf.AllowCredentials = true
	cf.AddAllowHeaders("X-Operator-Code", "X-Web-Code", "Authorization", "Credential")
	return cors.New(cf)
}
