package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sistema-notas/estoque/internal/core/domain"
	"sistema-notas/estoque/internal/core/ports"
)

type ProdutoHandler struct {
	repo ports.ProdutoRepository
}

// NewProdutoHandler
func NewProdutoHandler(repo ports.ProdutoRepository) *ProdutoHandler {
	return &ProdutoHandler{repo: repo}
}

// Criar
func (h *ProdutoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var p domain.Produto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Salvar(p)
	if err != nil {
		fmt.Println("Falha no banco de dados", err)
		http.Error(w, "Erro interno ao salvar o produto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensagem": "Produto salvo com sucesso", "id": %d}`, id)
}
