package init_app

import (
	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-label-pods/routes"
)

func InitGin(ginMode string) *gin.Engine {
	gin.SetMode(ginMode)

	r := gin.New()

	routes.InitLogger(r)

	routes.InitRouter(r)

	return r
}
