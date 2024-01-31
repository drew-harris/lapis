package views

import (
	"context"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func SendView(c *fiber.Ctx, cmp templ.Component) error {
	c.Set("Content-Type", "text/html")
	return cmp.Render(context.Background(), c.Response().BodyWriter())
}
