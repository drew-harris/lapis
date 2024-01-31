package standards

import (
	"fmt"
	"log"

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
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
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

		var msg string = "<div id=\"loglist\">test</div>"

		err := c2.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Println("got conn error")
			fmt.Println(err)
			return err
		}

	}
	return nil
}
