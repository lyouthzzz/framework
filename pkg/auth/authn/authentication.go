package authn

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lyouthzzz/framework/pkg/auth/user"
	"net/http"
)

const UserInfo = "authenticator:user"

// AuthN 系统主要用于认证 （Authentication），决定谁访问了系统。
// 根据配置验证用户凭证，通常是 Bearer Token 或者 Basic Auth Password。
// 解析用户信息，包括 Username，Groups，Uid 等。这些信息将用于之后的 Authz 系统授权。

type TokenAuthentication struct {
	Token string
}

type Authenticator interface {
	// 获取鉴权信息
	GetAuthentication(ctx context.Context, req *http.Request) (interface{}, error)
	// 增加鉴权信息
	AddAuthentication(ctx context.Context, req *http.Request, authentication interface{}) error
	// 写入鉴权信息
	WriteAuthentication(ctx context.Context, authentication interface{}, userInfo user.Info) error
	// 删除鉴权信息
	DeleteAuthentication(ctx context.Context, authentication interface{}) error
	// 鉴权验证
	Authenticate(ctx context.Context, authentication interface{}) (user.Info, error)
	// 鉴权失败
	AuthenticateFailedCB(w http.ResponseWriter, err error)
}

func NewAuthenticator(authN Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		authentication, err := authN.GetAuthentication(ctx, c.Request)
		if err != nil {
			authN.AuthenticateFailedCB(c.Writer, err)
			return
		}
		userInfo, err := authN.Authenticate(ctx, authentication)
		if err != nil {
			authN.AuthenticateFailedCB(c.Writer, err)
			return
		}

		c.Set(UserInfo, userInfo)
		c.Next()
	}
}
