package middleware

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil/resp"
	"slices"
)

type RoleRestrictionConfig struct {
	ClaimsKey string // Default: claims
	roles     []users.Role
}

func NewRoleRestrictionMiddleware(config *RoleRestrictionConfig, roles []users.Role) RoleRestrictionMiddleware {
	if config == nil {
		config = &RoleRestrictionConfig{}
	}

	if roles == nil || len(roles) == 0 {
		panic("zeroed")
	}

	config.roles = roles
	util.SetDefaultString(&config.ClaimsKey, "claims")

	return RoleRestrictionMiddleware{config: config}
}

// RoleRestrictionMiddleware will only look up on config.ClaimsKey which will be set by other middlewares,
// in this context will be handled by authorization middleware
type RoleRestrictionMiddleware struct {
	config *RoleRestrictionConfig
}

func (r RoleRestrictionMiddleware) Handle(ctx *gin.Context) {
	// Open claims
	claims, err := util.GetContextValue[*common.AccessTokenClaims](ctx, r.config.ClaimsKey)
	if err != nil {
		resp.Error(ctx, status.Error(status.AUTH_UNAUTHORIZED))
		ctx.Abort()
		return
	}

	role, err := users.NewRole(claims.Role)
	if err != nil {
		resp.Error(ctx, status.Error(status.JWT_TOKEN_MALFORMED))
		ctx.Abort()
		return
	}

	if !r.hasPermission(role) {
		resp.Error(ctx, status.Error(status.AUTH_UNAUTHORIZED))
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (r RoleRestrictionMiddleware) hasPermission(needle users.Role) bool {
	return slices.Contains(r.config.roles, needle)
}
