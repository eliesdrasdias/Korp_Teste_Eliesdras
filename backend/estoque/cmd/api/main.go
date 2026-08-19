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

// corsMiddleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Configurações de CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Configuração do banco
	connStr := "host=localhost port=5432 user=root password=rootpassword dbname=sistema_notas sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro fatal ao configurar banco:", err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Erro ao conectar!", err)
	}
	fmt.Println("Conetado ao PostgreSQL com sucesso!")

	// Injeção de dependência
	produtoRepo := repository.NewProdutoPostgres(db)
	produtoHandler := handler.NewProdutoHandler(produtoRepo)
	// Rotas
	http.HandleFunc("/produtos", corsMiddleware(produtoHandler.Criar))
	http.HandleFunc("/produtos/listar", corsMiddleware(produtoHandler.Listar))
	// Servidor
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}
