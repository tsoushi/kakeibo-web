package main

import (
	"kakeibo-web-server/handler/graph"
	"kakeibo-web-server/handler/graph/resolver"
	"kakeibo-web-server/handler/middleware"
	"kakeibo-web-server/repository"
	"kakeibo-web-server/usecase"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/gocraft/dbr/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/go-sql-driver/mysql"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mysqlConfig := mysql.Config{
		DBName:    os.Getenv("MYSQL_DATABASE"),
		User:      os.Getenv("MYSQL_USERNAME"),
		Passwd:    os.Getenv("MYSQL_USERPASS"),
		Addr:      os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"),
		Net:       "tcp",
		ParseTime: true,
	}

	dbrConn, err := dbr.Open("mysql", mysqlConfig.FormatDSN(), nil)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbrConn.SetMaxOpenConns(10)

	sess := dbrConn.NewSession(nil)

	repository := repository.NewRepository(sess)
	usecase := usecase.NewUsecase(repository)

	r := chi.NewRouter()
	graphQLRouter := chi.NewRouter()
	graphQLRouter.Use(middleware.MakeDebugAuth(repository))

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver.NewResolver(usecase),
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	graphQLRouter.Handle("/query", srv)
	r.Handle("/query", graphQLRouter)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
