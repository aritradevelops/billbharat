package authn

import (
	"fmt"
	"strings"

	"github.com/aritradevelops/billbharat/backend/product/internal/core/jwtutil"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/gofiber/fiber/v2"
)

const authUserKey = "auth_user"

func Middleware(jwtManager *jwtutil.JwtManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		accessToken := strings.TrimPrefix(bearer, "Bearer ")
		if accessToken == "" {
			accessToken = c.Cookies("access_token")
			logger.Info().Msg("Access token not found in header")
		}
		if accessToken == "" {
			logger.Info().Msg("Access token not found in cookies")
			return fiber.ErrUnauthorized
		}
		payload, err := jwtManager.Verify(accessToken)
		logger.Error().Err(err).Msg("Access token verification failed")
		if err != nil {
			return fiber.ErrUnauthorized
		}
		logger.Info().Msg("Access token verified successfully")
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
