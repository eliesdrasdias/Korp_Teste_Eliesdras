package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sistema-notas/estoque/internal/adapters/handler"
	"sistema-notas/estoque/internal/adapters/repository"

	_ "github.com/lib/pq"
)

func main() {
	// Configuração do banco
	connStr := "host=localhost port=5432 user=root password=rootpassword dbname=sistema_notas sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro fatal ao configurar banco:", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar! O banco pode estar desligado:", err)
	}

	fmt.Println("Conetado ao PostgreSQL com sucesso!")

	// Injeção de dependência
	produtoRepo := repository.NewProdutoPostgres(db)
	produtoHandler := handler.NewProdutoHandler(produtoRepo)
	// Rotas
	http.HandleFunc("/produtos", produtoHandler.Criar)
	// Servidor
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}
