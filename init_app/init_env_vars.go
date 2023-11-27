package init_app

import (
	"os"
)

func InitEnvVars() string {
	ginMode := "debug"
	ginModeVal, ginModePresent := os.LookupEnv("GIN_MODE")
	if ginModeVal == "" || !ginModePresent {
		os.Setenv("GIN_MODE", ginMode)
		ginModeVal = ginMode
	}

	return ginModeVal
}
