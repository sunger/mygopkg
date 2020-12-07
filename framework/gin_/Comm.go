import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 使gin具有跨域能力
func GinCors(r *gin.Engine) {

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"token", "content-type", "authorization", "origin", "accept", "x-requested-with"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}
	r.Use(cors.New(config))
	// cors；注意：生产环境可以注释掉
}
