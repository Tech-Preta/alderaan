package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Config contém a configuração do banco de dados
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresConnection cria uma nova conexão com o PostgreSQL
func NewPostgresConnection(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)                 // Máximo de conexões abertas
	db.SetMaxIdleConns(5)                  // Máximo de conexões ociosas
	db.SetConnMaxLifetime(5 * time.Minute) // Tempo de vida máximo de uma conexão

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar no banco: %w", err)
	}

	return db, nil
}

// Close fecha a conexão com o banco
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
