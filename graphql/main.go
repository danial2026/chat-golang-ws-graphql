package main

import (
	"chat-golang-react/chat/graphql/auth"
	"chat-golang-react/chat/graphql/graphql"
	"chat-golang-react/chat/graphql/resources"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

const defaultPort = "4000"

func main() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln("Error loading .env file : ", err)
	}

	port := os.Getenv("GRAPHQLPORT")
	if port == "" {
		port = defaultPort
	}

	res, err := resources.ConstructResource()
	if err != nil {
		log.Println("Error during resource construction:", err)

		return
	}

	resolver := graphql.Resolver{
		ROOMSDB: res.ROOMSDB,
		USERDB:  res.USERDB,
		Timeout: res.Timeout,
	}

	router := chi.NewRouter()

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"*",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           600, // Maximum value not ignored by any of major browsers
	}))

	// auth middleware
	router.Use(auth.Middleware())

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &resolver}))

	// TODO : coment me before pushing to production
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))

	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
