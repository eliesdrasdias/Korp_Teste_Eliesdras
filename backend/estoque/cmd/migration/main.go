package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	connStr := "host=localhost port=5432 user=root password=rootpassword dbname=sistema_notas sslmode=disable"

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
