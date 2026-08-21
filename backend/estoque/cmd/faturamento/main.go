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
	connStr := "host=localhost port=5432 user=root password=rootpassword dbname=sistema_notas sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro fatal ao configurar banco:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Erro ao conectar!", err)
	}
	fmt.Println("Faturamento conectado ao PostgreSQL com sucesso!")

	notaRepo := repository.NewNotaRepositoryPostgres(db)
	notaHandler := handler.NewNotaHandler(notaRepo)

	http.HandleFunc("/notas", corsMiddleware(notaHandler.Emitir))
	http.HandleFunc("/notas/imprimir", corsMiddleware(notaHandler.Imprimir))

	fmt.Println("Serviço de Faturamento rodando na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
