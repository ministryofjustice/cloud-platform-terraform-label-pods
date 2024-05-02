package routes

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	n "github.com/ministryofjustice/cloud-platform-label-pods/pkg/get_team"
	m "github.com/ministryofjustice/cloud-platform-label-pods/pkg/mutate"
	"github.com/ministryofjustice/cloud-platform-label-pods/utils"
)

func initMutatePod(r *gin.Engine) {
	r.POST("/mutate/pod", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()

		if err != nil {
			errObj := utils.Response{
				Status: http.StatusInternalServerError,
				Data:   nil,
			}

			utils.SendResponse(c, errObj)
		}

		getGithubTeamnameFn := n.InitGetGithubTeamName(n.GetTeamName)

		mutated, err := m.Mutate(body, getGithubTeamnameFn)
		if err != nil {
			errObj := utils.Response{
				Status: http.StatusInternalServerError,
				Data:   nil,
			}

			utils.SendResponse(c, errObj)
		}

		_, writeErr := c.Writer.Write(mutated)
		if writeErr != nil {
			errObj := utils.Response{
				Status: http.StatusInternalServerError,
				Data:   nil,
			}

			utils.SendResponse(c, errObj)
		}
	})
}
