package authn

import (
	"fmt"
	"strings"

	"github.com/aritradeveops/billbharat/backend/auth/internal/core/jwtutil"
	"github.com/gofiber/fiber/v2"
)

const authUserKey = "auth_user"

func Middleware(jwtManager *jwtutil.JwtManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		if bearer == "" {
			bearer = c.Cookies("access_token")
		}
		if bearer == "" {
			return fiber.ErrUnauthorized
		}
		accessToken := strings.TrimPrefix(bearer, "Bearer ")
		payload, err := jwtManager.Verify(accessToken)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals(authUserKey, payload)
		return c.Next()
	}
}

func GetUserFromContext(c *fiber.Ctx) (*jwtutil.JwtPayload, error) {
	userIn := c.Locals(authUserKey)
	if userIn == nil {
		return nil, fmt.Errorf("AuthenticatedUser is only available for protected routes")
	}
	payload, ok := userIn.(*jwtutil.JwtPayload)
	if !ok {
		return nil, fmt.Errorf("AuthenticatedUser is only available for protected routes")
	}
	return payload, nil
}
