package main

import (
	"database/sql"
	"fmt"
	"log"
	"sistema-notas/estoque/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	connStr := config.Get("DATABASE_URL", "host=localhost port=5432 user=postgres password=postgres dbname=sistema_notas sslmode=disable")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro fatal ao configurar banco:", err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Erro ao conectar!", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("Erro ao configurar o dialecto:", err)
	}

	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Fatal("Erro ao executar as migrations:", err)
	}

	fmt.Println("Migrations executadas com sucesso!")
}
