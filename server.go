package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/drew-harris/lapis/code"
	"github.com/drew-harris/lapis/graph"
	"github.com/drew-harris/lapis/graph/model"
	"github.com/drew-harris/lapis/players"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/posthog/posthog-go"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const defaultPort = "8080"

type CreatePlayerInput struct {
	Name string `json:"name"`
}

func main() {
	//Load env safely
	err := godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	engine := html.New("./views/", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())

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

	resolver := graph.NewResolver(db, pClient)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

	app.Get("/codes", func(c *fiber.Ctx) error {
		players, err := players.GetAllPlayers(db)
		if err != nil {
			return err
		}
		return c.Render("code", fiber.Map{
			"players": players,
		})
	})

	app.Get("/styles.css", func(c *fiber.Ctx) error {
		return c.SendFile("./views/styles.css")
	})

	app.Post("/hx/setup", func(c *fiber.Ctx) error {
		playerInput := CreatePlayerInput{}
		if err := c.BodyParser(&playerInput); err != nil {
			return err
		}

		player, err := code.RegisterPlayerWithNewCode(playerInput.Name, db)
		if err != nil {
			return err
		}

		pClient.Enqueue(posthog.Alias{
			DistinctId: player.ID,
			Alias:      player.Name,
			Timestamp:  time.Now(),
		})

		pClient.Enqueue(posthog.Identify{
			DistinctId: player.ID,
			Properties: posthog.NewProperties().Set("name", player.Name).Set("created", time.Now()),
			Timestamp:  time.Now(),
		})

		return c.Render("result", fiber.Map{
			"ID":   player.ID,
			"Name": player.Name,
		})
	})

	// Handle fiber with gql
	app.Use("/query", func(c *fiber.Ctx) error {
		test := adaptor.HTTPHandler(srv)
		return test(c)
	})

	// Handle fiber with gql
	app.Use("/", func(c *fiber.Ctx) error {
		test := adaptor.HTTPHandler(playground.Handler("Graphql playground", "/query"))
		return test(c)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	app.Listen(":" + port)
}
