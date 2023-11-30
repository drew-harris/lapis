package main

import (
	"fmt"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/drew-harris/lapis/graph"
	"github.com/drew-harris/lapis/graph/model"
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

	app := fiber.New()

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

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
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
