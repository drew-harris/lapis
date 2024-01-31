package standards

import (
	"context"
	"fmt"

	"github.com/drew-harris/lapis/graph/model"
	"github.com/drew-harris/lapis/views"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type StandardsController struct {
	db      *gorm.DB
	wsConns []*websocket.Conn
}

func New(db *gorm.DB) *StandardsController {
	return &StandardsController{db: db}
}

func (s *StandardsController) ShowStandardsPage(c *fiber.Ctx) error {
	return views.SendView(c, views.StandardsPage())
}

func (s *StandardsController) AddWebSocketConnection(c *websocket.Conn) error {
	s.wsConns = append(s.wsConns, c)
	fmt.Println("Added websocket connection")

	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			break
		}

		if err = c.WriteMessage(mt, msg); err != nil {
			break
		}
	}

	// Remove connection
	for i, c2 := range s.wsConns {
		if c2 == c {
			s.wsConns = append(s.wsConns[:i], s.wsConns[i+1:]...)
			break
		}
	}

	return nil
}

func (s *StandardsController) AddLog(log *model.Log) error {
	fmt.Println("Adding log to standars")
	for _, c2 := range s.wsConns {
		if c2 == nil {
			continue
		}
		if c2.Conn == nil {
			continue
		}

		writer, err := c2.NextWriter(websocket.TextMessage)
		if err != nil {
			return err
		}

		err = views.LogRow(*log).Render(context.Background(), writer)
		if err != nil {
			return err
		}
	}
	return nil
}
