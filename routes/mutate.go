package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/ministryofjustice/cloud-platform-label-pods/pkg/mutate"
	n "github.com/ministryofjustice/cloud-platform-label-pods/pkg/namespace"
	"github.com/ministryofjustice/cloud-platform-label-pods/utils"
)

func initMutatePod(r *gin.Engine) {
	r.POST("/mutate/pod", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()

		if err != nil {
			errObj := utils.Response{
				Status: http.StatusInternalServerError,
				Data:   nil,
			}

			utils.SendResponse(c, errObj)
		}

		getGithubTeamnameFn := n.InitGetGithubTeamName(n.GetTeamNameFromNs)

		mutated, err := m.Mutate(body, getGithubTeamnameFn)
		if err != nil {
			errObj := utils.Response{
				Status: http.StatusInternalServerError,
				Data:   nil,
			}

			utils.SendResponse(c, errObj)
		}

		c.Writer.Write(mutated)
	})
}
