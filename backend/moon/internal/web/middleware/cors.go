package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "x-jwt-token, x-refresh-token")

		if ctx.Request.Method == "OPTIONS" {
			log.Printf("CORS: OPTIONS request to %s from %s", ctx.Request.URL.Path, origin)
			ctx.AbortWithStatus(204)
			return
		}

		log.Printf("CORS: %s request to %s from %s", ctx.Request.Method, ctx.Request.URL.Path, origin)
		ctx.Next()
	}
}
