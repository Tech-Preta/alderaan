package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/williamkoller/golang-domain-driven-design/docs"
	product_repository "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/repository"
	product_handlers "github.com/williamkoller/golang-domain-driven-design/internal/infra/http/handlers"
	product_router "github.com/williamkoller/golang-domain-driven-design/internal/infra/http/router"
	"github.com/williamkoller/golang-domain-driven-design/internal/infra/persistence"
	"github.com/williamkoller/golang-domain-driven-design/internal/metrics"
	"github.com/williamkoller/golang-domain-driven-design/internal/shared/config"
	"github.com/williamkoller/golang-domain-driven-design/internal/shared/database"
	shared_events "github.com/williamkoller/golang-domain-driven-design/internal/shared/domain/events"
)

//	@title			Servidor HTTP com Domain Driven Design
//	@version		1.0
//	@description	API RESTful com DDD, Gin e Graceful Shutdown
//	@contact.name	API Support
//	@contact.url	https://github.com/williamkoller
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//	@host			localhost:8080
//	@BasePath		/api/v1

func main() {
	// Carregar configurações
	cfg := config.Load()

	// Conectar ao banco de dados
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco de dados: %v", err)
	}
	defer database.Close(db)

	log.Println("✅ Conectado ao banco de dados PostgreSQL")

	// Inicializar componentes
	dispatcher := shared_events.NewEventDispatcher()

	// Usar repositório PostgreSQL ao invés de in-memory
	var repo product_repository.IProductRepository
	if db != nil {
		repo = persistence.NewPostgresProductRepository(db)
		log.Println("📊 Usando repositório PostgreSQL")
	} else {
		repo = product_repository.NewRepository()
		log.Println("💾 Usando repositório in-memory")
	}

	m := metrics.NewMetrics()

	productHandler := product_handlers.NewProductHandler(repo, dispatcher, m)

	r := product_router.SetupProductRouter(productHandler, m)

	server := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           r,
		ReadHeaderTimeout: 30 * time.Second,
	}

	go func() {
		fmt.Printf("\n🚀 Server running on http://localhost:%s\n", cfg.Server.Port)
		fmt.Printf("📊 Metrics available at http://localhost:%s/metrics\n", cfg.Server.Port)
		fmt.Printf("📚 Swagger UI at http://localhost:%s/swagger/index.html\n", cfg.Server.Port)
		fmt.Printf("🗄️  Database: %s@%s:%d/%s\n\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v\n", err)
		}
	}()

	GracefulShutdown(server, db, 5*time.Second)
}

func GracefulShutdown(server *http.Server, db interface{ Close() error }, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("\n⚠️  Received termination signal. Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Error during server shutdown: %v\n", err)
	} else {
		log.Println("✅ HTTP server shut down gracefully")
	}

	// Close database connection
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("❌ Error closing database: %v\n", err)
		} else {
			log.Println("✅ Database connection closed")
		}
	}

	log.Println("👋 Shutdown complete")
}
