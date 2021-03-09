package reponseHandler

import "github.com/gofiber/fiber/v2"

type reponseMessage struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status,omitempty"`
}

func ReponseMsg(c *fiber.Ctx, code int, status string, msg string, data interface{}) error {
	reponseMessage := &reponseMessage{
		Code:    code,
		Message: msg,
		Data:    data,
		Status:  status,
	}

	return c.Status(code).JSON(reponseMessage)
}
