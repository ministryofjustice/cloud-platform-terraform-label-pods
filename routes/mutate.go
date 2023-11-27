package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/ministryofjustice/cloud-platform-label-pods/pkg/mutate"
	"github.com/ministryofjustice/cloud-platform-label-pods/utils"
)

func InitGetCheckPR(r *gin.Engine) {
	r.POST("/mutate/pod", func(c *gin.Context) {
		// if not a system nameespace then continue

		body, err := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()

		if err != nil {
			errObj := utils.Response{
				Status: http.StatusInternalServerError,
				Data:   nil,
			}

			utils.SendResponse(c, errObj)
		}

		mutated, err := m.Mutate(body, "random")
		if err != nil {
			errObj := utils.Response{
				Status: http.StatusInternalServerError,
				Data:   nil,
			}

			utils.SendResponse(c, errObj)
		}

		// TODO: this needs to be sent back in the correct obj adm review
		obj := utils.Response{
			Status: http.StatusOK,
			Data:   mutated,
		}
		utils.SendResponse(c, obj)
	})
}
