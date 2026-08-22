package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sistema-notas/estoque/internal/adapters/handler"
	"sistema-notas/estoque/internal/adapters/repository"
	"sistema-notas/estoque/internal/config"

	_ "github.com/lib/pq"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	connStr := config.Get("DATABASE_URL", "host=localhost port=5432 user=postgres password=postgres dbname=sistema_notas sslmode=disable")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro fatal ao configurar banco:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Erro ao conectar!", err)
	}
	fmt.Println("Estoque conectado ao PostgreSQL com sucesso!")

	produtoRepo := repository.NewProdutoPostgres(db)
	produtoHandler := handler.NewProdutoHandler(produtoRepo)

	http.HandleFunc("/produtos", corsMiddleware(produtoHandler.Criar))
	http.HandleFunc("/produtos/listar", corsMiddleware(produtoHandler.Listar))
	http.HandleFunc("/produtos/baixa", corsMiddleware(produtoHandler.BaixarEstoque))

	fmt.Println("Serviço de Estoque rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
