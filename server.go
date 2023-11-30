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
	"github.com/gofiber/template/html/v2"

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
		return
	}

	resolver := graph.NewResolver(db)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

	app.Get("/codes", func(c *fiber.Ctx) error {
		return c.Render("code", nil)
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
		return c.Render("result", fiber.Map{
			"code": player.ID,
			"name": player.Name,
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
