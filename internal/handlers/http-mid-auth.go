package handlers

import (
	"crypto/rsa"

	"github.com/HangoKub/Hango-service/pkg/reponseHandler"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

type middlewareAuth struct {
	privKey *rsa.PrivateKey
}

type MiddlewareAuth interface {
	Prodected() func(c *fiber.Ctx) error
}

func NewMiddlewareAuth(privKey *rsa.PrivateKey) MiddlewareAuth {
	return &middlewareAuth{
		privKey,
	}
}

func (mid *middlewareAuth) Prodected() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:    mid.privKey.Public(),
		SigningMethod: "RS512",
		ErrorHandler:  jwtError,
		TokenLookup:   "header:Authorization",
		AuthScheme:    "Bearer",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return reponseHandler.ReponseMsg(c, fiber.StatusUnauthorized, "failed", "Missing or malformed JWT", nil)
	}
	return reponseHandler.ReponseMsg(c, fiber.StatusUnauthorized, "failed", "Invalid or expired JWT", nil)
}
