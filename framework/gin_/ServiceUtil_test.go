package gin_

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBootStrapGroup(t *testing.T) {

	r := gin.New()
	//gin_.BootStrap(r)
	v1 := r.Group("v1")

	BootStrapGroup(v1)
}
