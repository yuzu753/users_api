package runtime

import (
	"users_api/src/interface/web/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(uc *controller.UserController) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/:tenant_id")
	{
		v1.GET("/Users", uc.Search)
	}
	
	return r
}