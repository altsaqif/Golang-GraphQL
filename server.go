package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/altsaqif/graphql-go/cmd/configs"
	"github.com/altsaqif/graphql-go/cmd/middlewares"
	"github.com/altsaqif/graphql-go/graph"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize resolver with database instance
	configs.ConnectionDB()

	// Initialize Gorilla Mux router
	router := chi.NewRouter()

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Route does not exist!"))
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("Method is not valid!"))
	})

	// Route for GraphQL Playground
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))

	router.Route("/api", func(auth chi.Router) {

		resolver := graph.Resolver{DB: configs.DB}

		// Initialize GraphQL server
		srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

		// Initialize session middleware
		auth.Use(middlewares.Middleware)

		// Route for GraphQL endpoint
		auth.Handle("/query", srv)

	})

	// Start HTTP server
	log.Println("Server is running on http://localhost:9600")
	log.Fatal(http.ListenAndServe(":9600", router))
}
