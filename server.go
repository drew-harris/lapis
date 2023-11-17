package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/drew-harris/lapis/graph"
	"github.com/drew-harris/lapis/graph/model"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		db.AutoMigrate(&model.Player{})
		db.AutoMigrate(&model.Save{})
		db.AutoMigrate(&model.Log{})
		db.AutoMigrate(&model.CustomNode{})
		return
	}

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{}).Handler)

	resolver := graph.NewResolver(db)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
