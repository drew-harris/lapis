package realtime

import (
	"context"
	"fmt"

	"github.com/gofiber/contrib/websocket"

	"github.com/a-h/templ"
)

type Service struct {
	wsConns []*websocket.Conn
}

func New() *Service {
	return &Service{}
}

func (s *Service) AddWebSocketConnection(c *websocket.Conn) error {
	s.wsConns = append(s.wsConns, c)

	var (
		mt  int
		msg []byte
		err error
	)

	// Say hello
	if err = c.WriteMessage(websocket.TextMessage, []byte("slif")); err != nil {
		fmt.Println(err)
		return err
	}

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

func (s *Service) SendHtml(cmp templ.Component) error {
	for _, c2 := range s.wsConns {
		if c2 == nil {
			continue
		}

		writer, err := c2.NextWriter(websocket.TextMessage)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = cmp.Render(context.Background(), writer)

		writer.Close()

		fmt.Println("Actually sent")
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}
