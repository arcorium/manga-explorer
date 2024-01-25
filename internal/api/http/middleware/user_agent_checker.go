package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil/resp"
)

type UserAgentConfig struct {
	Key            string // Default: user_agent
	allowArbitrary bool
}

func NewUserAgentChecker(allowArbitrary bool, config UserAgentConfig) UserAgentCheckerMiddleware {
	util.SetDefaultString(&config.Key, "user_agent")
	config.allowArbitrary = allowArbitrary
	return UserAgentCheckerMiddleware{Config: &config}
}

type UserAgentCheckerMiddleware struct {
	Config *UserAgentConfig
}

func (u UserAgentCheckerMiddleware) Handle(ctx *gin.Context) {
	deviceName := ""
	str := ctx.GetHeader("User-Agent")
	if len(str) == 0 {
		if u.Config.allowArbitrary {
			deviceName = "Unknown"
		} else {
			resp.Error(ctx, status.Error(status.USER_AGENT_UNKNOWN_ERROR))
			ctx.Abort()
			return
		}
	}
	ua := useragent.Parse(str)
	deviceName = users.ParseDeviceName(&ua)

	ctx.Set(u.Config.Key, &users.Device{Name: deviceName})
}
