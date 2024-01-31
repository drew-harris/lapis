package main

import (
	"fmt"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/drew-harris/lapis/code"
	"github.com/drew-harris/lapis/graph"
	"github.com/drew-harris/lapis/graph/model"
	"github.com/drew-harris/lapis/standards"
	"github.com/gofiber/contrib/websocket"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/posthog/posthog-go"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	//Load env safely
	err := godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := fiber.New(fiber.Config{})
	app.Use(cors.New())
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = fmt.Errorf("failed to connect database")
		panic(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		db.AutoMigrate(&model.Player{})
		db.AutoMigrate(&model.Save{})
		db.AutoMigrate(&model.Log{})
		db.AutoMigrate(&model.CustomNode{})
		db.AutoMigrate(&model.Position{})

		return
	}

	// Create Posthog client
	posthogKey := os.Getenv("POSTHOG_KEY")
	if posthogKey == "" {
		_ = fmt.Errorf("failed to find posthog key")
		panic(err)
	}
	pClient, _ := posthog.NewWithConfig(posthogKey, posthog.Config{
		Endpoint: "https://app.posthog.com",
	})

	standardsController := standards.New(db)

	resolver := graph.NewResolver(db, pClient, standardsController)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

	app.Get("/styles.css", func(c *fiber.Ctx) error {
		return c.SendFile("./dist/out.css")
	})

	// Handle fiber with gql
	app.Use("/query", func(c *fiber.Ctx) error {
		test := adaptor.HTTPHandler(srv)
		return test(c)
	})

	codeHandler := code.CreateCodeHandler(db, pClient)
	app.Get("/codes", codeHandler.GetCodes)
	app.Post("/hx/setup", codeHandler.SetupPlayer)

	app.Get("/standards", standardsController.ShowStandardsPage)

	// Live log view
	app.Get("/ws/logs", websocket.New(func(c *websocket.Conn) {
		standardsController.AddWebSocketConnection(c)
	}, websocket.Config{}))

	// Handle fiber with gql
	app.Use("/", func(c *fiber.Ctx) error {
		test := adaptor.HTTPHandler(playground.Handler("Graphql playground", "/query"))
		return test(c)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	app.Listen(":" + port)
}
