package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"

	"github.com/ybotet/pz8-pipelineCICD-go/services/graphql/graph"
	"github.com/ybotet/pz8-pipelineCICD-go/services/graphql/graph/generated"
	"github.com/ybotet/pz8-pipelineCICD-go/services/graphql/internal/repository"
	"github.com/ybotet/pz8-pipelineCICD-go/shared/logger"
)

func main() {
	// Usar tu logger existente
	log := logger.New(logger.Config{
		ServiceName: "graphql",
		Environment: "development",
		LogLevel:    "info",
		JSONFormat:  false,
	})
	
	log.Info("Starting GraphQL service...")

	port := getEnv("GRAPHQL_PORT", "8090")
	
	// Configuración de BD
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "tasksuser")
	dbPassword := getEnv("DB_PASSWORD", "taskspass")
	dbName := getEnv("DB_NAME", "tasksdb")
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()
	
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	
	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("Database ping failed")
	}
	log.Info("Connected to PostgreSQL")

	// Repositorio con PostgreSQL
	taskRepo := repository.NewPostgresTaskRepository(db, log)
	
	// Resolver
	resolver := graph.NewResolver(taskRepo, log)
	
	// Servidor GraphQL
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	
	http.Handle("/", playground.Handler("GraphQL Playground - Tasks API", "/query"))
	http.Handle("/query", srv)

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof("GraphQL server listening on port %s", port)
		log.Infof("Playground available at http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down GraphQL server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server shutdown error")
	}
	
	db.Close()
	log.Info("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}