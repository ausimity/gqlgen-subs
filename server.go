package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/websocket"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ausimity/gqlgen-todos/graph"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// add ws transport configured by ourselves
	srv.AddTransport(&transport.Websocket{
    		Upgrader: websocket.Upgrader{
        	//ReadBufferSize:  1024,
        	//WriteBufferSize: 1024,
        	CheckOrigin: func(r *http.Request) bool {
            		// add checking origin logic to decide return true or false
            		return true
        	},
    		},
    		KeepAlivePingInterval: 10 * time.Second,
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
