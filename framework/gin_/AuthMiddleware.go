package gin_

import (
	"strings"

	"github.com/gin-gonic/gin"
	mycs "github.com/sunger/mygopkg/framework/casbin"
)

//无域认证
func AuthMiddleware() Middleware {
	return func(next Endpoint) Endpoint {
		return func(c *gin.Context, request interface{}) (response interface{}, err error) {
			if c.Request.Header.Get("token") == "" {
				c.AbortWithStatusJSON(400, gin.H{"message": "token required"})
				return
			}
			token := c.Request.Header.Get("token")
			user, err := ParseToken(token)
			if err != nil {
				c.AbortWithStatusJSON(403, gin.H{"message": err})
				return
			}
			access, err := mycs.E.Enforce(user.ID, c.Request.RequestURI, c.Request.Method)
			if err != nil || !access {
				c.AbortWithStatusJSON(403, gin.H{"message": "forbidden"})
				return
			}
			return next(c, request)
		}
	}
}

//有域认证
func AuthWithDomainMiddleware() Middleware {
	return func(next Endpoint) Endpoint {
		return func(c *gin.Context, request interface{}) (response interface{}, err error) {
			if c.Request.Header.Get("token") == "" {
				c.AbortWithStatusJSON(400, gin.H{"message": "token required"})
				return
			}

			token := c.Request.Header.Get("token")
			user, err := ParseToken(token)
			if err != nil {
				c.AbortWithStatusJSON(403, gin.H{"message": err})
				return
			}

			// user, _ := c.Get("id")
			domain := c.Request.Header.Get("Host")
			uri := strings.TrimPrefix(c.Request.RequestURI, "/"+domain) // /domain/depts => /depts
			access, err := mycs.E.Enforce(user.ID, domain, uri, c.Request.Method)
			if err != nil || !access {
				c.AbortWithStatusJSON(403, gin.H{"message": "forbidden"})
				return
			}
			return next(c, request)
		}
	}
}
